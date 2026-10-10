package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ateto3/gopherml/ensemble"
	"github.com/ateto3/gopherml/persistence"
)

func main() {
	fmt.Println("=== Model Persistence Example ===")

	// Datos de ejemplo: [longitud, ancho] → 0 o 1
	X := [][]float64{
		{1.0, 2.0}, {1.5, 1.8}, {1.2, 2.5}, {1.1, 1.9},
		{5.0, 6.0}, {5.5, 5.8}, {5.2, 6.5}, {4.9, 5.9},
	}
	y := []float64{0, 0, 0, 0, 1, 1, 1, 1}

	// 1. Entrenar
	rf := ensemble.NewRandomForestClassifier()
	rf.NEstimators = 30
	rf.MaxDepth = 4
	rf.RandomState = 42
	if err := rf.Fit(X, y); err != nil {
		log.Fatalf("Fit: %v", err)
	}
	original, err := rf.Predict(X)
	if err != nil {
		log.Fatalf("Predict: %v", err)
	}
	fmt.Println("Predicciones originales:", original)

	// 2. Guardar (formato gob)
	const path = "persistence_example.gob"
	if err := persistence.Save(path, rf); err != nil {
		log.Fatalf("Save: %v", err)
	}
	defer os.Remove(path)
	fmt.Println("Modelo guardado en", path)

	// 3. Cargar en un modelo NUEVO — la integridad se valida automáticamente
	loaded := ensemble.NewRandomForestClassifier()
	if err := persistence.Load(path, loaded); err != nil {
		log.Fatalf("Load: %v", err)
	}
	reloaded, err := loaded.Predict(X)
	if err != nil {
		log.Fatalf("Predict tras load: %v", err)
	}
	fmt.Println("Predicciones tras cargar:", reloaded)
	fmt.Println("¿Coinciden?", fmt.Sprint(original) == fmt.Sprint(reloaded))

	// 4. Detectar archivos corruptos (la validación falla con un error claro)
	if err := os.WriteFile(path, []byte("esto no es un gob válido"), 0o644); err != nil {
		log.Fatalf("WriteFile: %v", err)
	}
	corrupto := ensemble.NewRandomForestClassifier()
	if err := persistence.Load(path, corrupto); err != nil {
		fmt.Printf("Archivo corrupto detectado correctamente:\n  %v\n", err)
	} else {
		log.Fatal("ERROR: se esperaba un error al cargar un archivo corrupto")
	}

	fmt.Println("\n=== Example Complete ===")
}
