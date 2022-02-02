package handler

import (
	"net/http"
	"se_uf/gator_snapstore/models"

	"gorm.io/gorm"
)

func GetAllImages(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var allProductImages []models.ProductCatalogue
	var allImages []models.Image
	allImagesFromDB := DB.Find(&allImages)
	rows, err := allImagesFromDB.Rows()
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var image models.Image
	for rows.Next() {
		DB.ScanRows(rows, &image)
		genres, err := getGenresOfImage(DB, w, image.ImageId)
		if err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		productCatalogueImage := models.ProductCatalogue{
			ImageId:   image.ImageId,
			Price:     image.Price,
			Title:     image.Title,
			WImageURL: image.WImageURL,
			Genre:     genres,
		}
		allProductImages = append(allProductImages, productCatalogueImage)
	}
	SendJSONResponse(w, http.StatusOK, allProductImages)
}

func getGenresOfImage(DB *gorm.DB, w http.ResponseWriter, imageId int) ([]string, error) {
	var genres []string
	var genre []models.Genre
	// Another way of finding data:
	// genresOfCurrentImage := DB.Where("image_id = ?", strconv.Itoa(imageId)).Find(&genre)
	genresOfCurrentImage := DB.Where(&models.Genre{ImageId: imageId}).Find(&genre)
	genreRows, err := genresOfCurrentImage.Rows()
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return genres, err
	}
	defer genreRows.Close()
	var currentGenre models.Genre
	for genreRows.Next() {
		DB.ScanRows(genreRows, &currentGenre)
		genres = append(genres, currentGenre.GenreType)
	}
	return genres, nil
}

func GetGenreCategories(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var genreCategories []models.GenreCategories
	allCategoriesFromDB := DB.Find(&genreCategories)
	genreCategoriesRows, err := allCategoriesFromDB.Rows()
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer genreCategoriesRows.Close()
	var allGenreCategories []string
	var genreCategory models.GenreCategories
	for genreCategoriesRows.Next() {
		DB.ScanRows(genreCategoriesRows, &genreCategory)
		allGenreCategories = append(allGenreCategories, genreCategory.Category)
	}
	SendJSONResponse(w, http.StatusOK, allGenreCategories)
}
