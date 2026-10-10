package main

import (
	"fmt"
	"log"

	"github.com/ateto3/gopherml/preprocessing"
)

func main() {
	fmt.Println("=== Preprocessing Pipeline Example ===")

	X := [][]float64{
		{1500, 3, 10},
		{2000, 4, 5},
		{1200, 2, 15},
		{2500, 4, 2},
		{1800, 3, 8},
		{3000, 5, 1},
	}

	// 1. StandardScaler — estandarizar a media 0 / desviación 1
	scaler := preprocessing.NewStandardScaler()
	Xs, err := scaler.FitTransform(X)
	if err != nil {
		log.Fatalf("StandardScaler: %v", err)
	}
	fmt.Printf("StandardScaler: %d muestras, primera fila = %.3f %.3f %.3f\n",
		len(Xs), Xs[0][0], Xs[0][1], Xs[0][2])

	// Transformaciones futuras usan los MISMOS parámetros del Fit inicial
	again, err := scaler.Transform([][]float64{{1500, 3, 10}})
	if err != nil {
		log.Fatalf("Transform: %v", err)
	}
	fmt.Printf("Transform repetido (misma fila): %.3f %.3f %.3f\n",
		again[0][0], again[0][1], again[0][2])

	// Invertir la transformación
	original, err := scaler.InverseTransform(Xs)
	if err != nil {
		log.Fatalf("InverseTransform: %v", err)
	}
	fmt.Printf("InverseTransform: primera fila = %.0f %.0f %.0f\n",
		original[0][0], original[0][1], original[0][2])

	// 2. MinMaxScaler — escalar al rango [0, 1]
	minMax := preprocessing.NewMinMaxScaler()
	Xmm, err := minMax.FitTransform(X)
	if err != nil {
		log.Fatalf("MinMaxScaler: %v", err)
	}
	fmt.Printf("MinMaxScaler: primera fila = %.3f %.3f %.3f\n",
		Xmm[0][0], Xmm[0][1], Xmm[0][2])

	// 3. LabelEncoder — codificar etiquetas de texto a enteros
	le := preprocessing.NewLabelEncoder()
	y := []string{"negativo", "positivo", "negativo", "positivo"}
	le.Fit(y)
	encoded, err := le.Transform(y)
	if err != nil {
		log.Fatalf("LabelEncoder.Transform: %v", err)
	}
	fmt.Printf("LabelEncoder: %v → %v\n", y, encoded)

	decoded, err := le.InverseTransform(encoded)
	if err != nil {
		log.Fatalf("LabelEncoder.InverseTransform: %v", err)
	}
	fmt.Printf("InverseTransform: %v → %v\n", encoded, decoded)

	// 4. PolynomialFeatures — features polinómicas de grado 2
	pf := preprocessing.NewPolynomialFeatures(2)
	Xp, err := pf.FitTransform(Xs)
	if err != nil {
		log.Fatalf("PolynomialFeatures: %v", err)
	}
	fmt.Printf("PolynomialFeatures(2): %d features → %d features\n",
		len(Xs[0]), len(Xp[0]))

	// ⚠️ Transform antes de Fit devuelve un error explícito, nunca falla en silencio
	unfitted := preprocessing.NewStandardScaler()
	if _, err := unfitted.Transform(X); err != nil {
		fmt.Printf("Control de errores OK: %v\n", err)
	}

	fmt.Println("\n=== Example Complete ===")
}
