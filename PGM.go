package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type PGM struct {
	Data        [][]uint8
	Width       int
	Height      int
	MagicNumber string
	MaxValue    int
}

type PBM struct {
	Data [][]bool
}

func main() {
	filename := "example.pgm"
	pgm, err := ReadPGM(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Manipulation du PGM
	pgm.Invert()
	pgm.Flip()

	// Conversion en PBM
	pbm := pgm.ToPBM()

	// Sauvegarde du PBM
	pbm.Save("output.pbm")
}

func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture du numéro magique : %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, fmt.Errorf("numéro magique invalide : %s", magicNumber)
	}

	dimensions, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture des dimensions : %v", err)
	}
	var width, height int
	_, err = fmt.Sscanf(strings.TrimSpace(dimensions), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("dimensions invalides : %v", err)
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("dimensions invalides : la largeur et la hauteur doivent être positives")
	}

	maxValue, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture de la valeur maximale : %v", err)
	}
	maxValue = strings.TrimSpace(maxValue)
	var max int
	_, err = fmt.Sscanf(maxValue, "%d", &max)
	if err != nil {
		return nil, fmt.Errorf("valeur maximale invalide : %v", err)
	}

	data := make([][]uint8, height)
	expectedBytesPerPixel := 1

	if magicNumber == "P2" {
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("erreur de lecture des données à la ligne %d : %v", y, err)
			}
			fields := strings.Fields(line)
			rowData := make([]uint8, width)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index hors limite à la ligne %d", y)
				}
				var pixelValue uint8
				_, err := fmt.Sscanf(field, "%d", &pixelValue)
				if err != nil {
					return nil, fmt.Errorf("erreur de conversion de la valeur du pixel à la ligne %d, colonne %d : %v", y, x, err)
				}
				rowData[x] = pixelValue
			}
			data[y] = rowData
		}
	} else if magicNumber == "P5" {
		for y := 0; y < height; y++ {
			row := make([]byte, width*expectedBytesPerPixel)
			n, err := reader.Read(row)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("fin de fichier inattendue à la ligne %d", y)
				}
				return nil, fmt.Errorf("erreur de lecture des données de pixels à la ligne %d : %v", y, err)
			}
			if n < width*expectedBytesPerPixel {
				return nil, fmt.Errorf("fin de fichier inattendue à la ligne %d, attendu %d octets, obtenu %d", y, width*expectedBytesPerPixel, n)
			}

			rowData := make([]uint8, width)
			for x := 0; x < width; x++ {
				pixelValue := uint8(row[x*expectedBytesPerPixel])
				rowData[x] = pixelValue
			}
			data[y] = rowData
		}
	}

	return &PGM{Data: data, Width: width, Height: height, MagicNumber: magicNumber, MaxValue: max}, nil
}

func (pgm *PGM) Invert() {
	for i, row := range pgm.Data {
		for j, value := range row {
			pgm.Data[i][j] = uint8(pgm.MaxValue) - value
		}
	}
}

func (pgm *PGM) Flip() {
	for x := 0; x < pgm.Height; x++ {
		for i, j := 0, pgm.Width-1; i < j; i, j = i+1, j-1 {
			pgm.Data[x][i], pgm.Data[x][j] = pgm.Data[x][j], pgm.Data[x][i]
		}
	}
}

// Fonction ToPBM manquante dans le code initial, ajoutée ici
func (pgm *PGM) ToPBM() *PBM {
	// Logique de conversion de PGM à PBM
	return nil
}

// Fonction Save manquante dans le code initial, ajoutée ici
func (pbm *PBM) Save(filename string) error {
	// Logique de sauvegarde pour PBM
	return nil
}
