package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
}

// PBM représente une image PBM.
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM lit une image PBM à partir d'un fichier et renvoie une structure qui représente l'image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	magicNumber := scanner.Text()
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, errors.New("unsupported file type")
	}

	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, errors.New("invalid image dimensions")
	}

	width, _ := strconv.Atoi(dimensions[0])
	height, _ := strconv.Atoi(dimensions[1])

	var data [][]bool
	if magicNumber == "P1" {
		for scanner.Scan() {
			line := scanner.Text()
			row := make([]bool, width)
			for i, char := range strings.Fields(line) {
				pixel, _ := strconv.Atoi(char)
				row[i] = pixel == 1
			}
			data = append(data, row)
		}
	} else if magicNumber == "P4" {
		// Créer un buffer pour lire les données binaires
		reader := bufio.NewReader(file)
		// Ignorer les espaces blancs qui pourraient exister après les dimensions
		reader.Discard(width % 8)
		for y := 0; y < height; y++ {
			row := make([]bool, width)
			for x := 0; x < width; x += 8 {
				// Lire un octet (8 bits) à la fois
				b, err := reader.ReadByte()
				if err != nil {
					return nil, err
				}
				// Convertir l'octet en booléens
				for i := 0; i < 8; i++ {
					// Vérifier si le bit à la position i est défini
					row[x+i] = b&(1<<(7-i)) != 0
				}
			}
			data = append(data, row)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At renvoie la valeur du pixel en (x, y).
func (pbm *PBM) At(x, y int) bool {
	if len(pbm.data) == 0 || x < 0 || y < 0 || x >= pbm.width || y >= pbm.height {
		return false
	}
	return pbm.data[y][x]
}

// Set définit la valeur du pixel à (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	if y >= 0 && y < pbm.height && x >= 0 && x < pbm.width {
		pbm.data[y][x] = value
	}
}

// Save enregistre l'image PBM dans un fichier et renvoie une erreur en cas de problème.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écrire le nombre magique et les dimensions
	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Écrire les pixels
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pixel {
				fmt.Fprint(file, "1 ")
			} else {
				fmt.Fprint(file, "0 ")
			}
		}
		fmt.Fprintln(file)
	}

	return nil
}

// Inverser inverse les couleurs de l'image PBM.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip retourne l'image PBM horizontalement.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width/2; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}

// Flop floppe l'image PBM verticalement.
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

// SetMagicNumber définit le nombre magique de l'image PBM.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
