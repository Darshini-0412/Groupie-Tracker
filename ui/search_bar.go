package ui

import (
	"strconv"
	"strings"

	"groupie-tracker/models"
)

func RechercherArtistes(artistes []models.Artist, recherche string) []models.Artist {
	if recherche == "" {
		return artistes
	}

	rechercheLower := strings.ToLower(recherche)
	var result []models.Artist

	for _, a := range artistes {
		if strings.Contains(strings.ToLower(a.Name), rechercheLower) {
			result = append(result, a)
			continue
		}

		membreTrouve := false
		for _, membre := range a.Members {
			if strings.Contains(strings.ToLower(membre), rechercheLower) {
				membreTrouve = true
				break
			}
		}
		if membreTrouve {
			result = append(result, a)
			continue
		}

		if strings.Contains(strconv.Itoa(a.CreationDate), recherche) {
			result = append(result, a)
			continue
		}

		if strings.Contains(strings.ToLower(a.FirstAlbum), rechercheLower) {
			result = append(result, a)
			continue
		}
	}

	return result
}
