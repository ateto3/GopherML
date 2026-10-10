# 🤝 Contribuir a GopherML

¡Gracias por tu interés en mejorar GopherML! Este documento recoge el flujo de trabajo, los estándares de código y las directrices de rendimiento que seguimos. Al enviar una pull request aceptas también nuestro [Código de Conducta](CODE_OF_CONDUCT.md).

---

## 1. Puesta en marcha

```bash
# 1. Haz un fork del repositorio y clónalo
git clone https://github.com/TU-USUARIO/gopherml.git
cd gopherml

# 2. Verifica que todo compila y los tests pasan
go build ./...
go vet ./...
go test ./... -race
```

**Requisitos:** Go 1.23 o superior. El proyecto no tiene dependencias externas, así que no hace falta `go mod download` extra.

> **Nota:** `go-nlp/` es un módulo independiente con su propio `go.mod`. Si tocas ese código,
> ejecuta también `cd go-nlp && go build ./... && go test ./...`.

---

## 2. Estándares de código

1. **Go idiomático:** sin patrones de Python ni Java. Interfaces pequeñas, nombres cortos, cero abstracción innecesaria.
2. **Gestión de errores explícita:** nunca ignorar el error devuelto por `Fit`, `Predict`, `Load` ni `Save`. Los errores se envuelven con contexto (`fmt.Errorf("...: %w", err)`).
3. **Constructores:** usa siempre los constructores `New...()`; configuran hiperparámetros por defecto críticos (`NEstimators`, `MaxIter`, etc.).
4. **Formato:** el código debe pasar `gofmt -w .` y `go vet ./...` sin avisos.
5. **Documentación de exportados:** toda función, tipo o método exportado lleva un comentario de doc en inglés que empieza por su nombre (convención de Go).

---

## 3. Directrices de rendimiento (críticas)

Este **no** es un proyecto de juguete: buscamos eficiencia cercana a C++ para inferencia de baja latencia.

1. **Prioridad de rendimiento:** prefiere *slices* pre-alojados y evita asignaciones de memoria en el *hot path*. Usa `PredictSingle` en lugar de `Predict` dentro de bucles de alta frecuencia.
2. **Concurrencia nativa:** usa `sync.WaitGroup` y canales solo cuando el paralelismo compensa el coste de la comunicación. Optimizamos para CPUs multinúcleo (ARM64 y x86_64).
3. **Zero dependencies:** el core depende solo de la librería estándar de Go. No añadir dependencias externas ni Cgo.
4. **Seguridad de datos:** todos los datos de entrada pasan por `utils.CheckFinite` (o `CheckFiniteSingle` en el hot path) antes de llegar a cualquier modelo.
5. **Sin regresiones silenciosas:** si optimizas, incluye o actualiza un benchmark en `*_test.go` que lo demuestre.

---

## 4. Directrices para asistentes IA (Windsurf / Antigravity / Copilot…)

Si estás generando código con una IA para este repositorio, sigue estas reglas sin excepción:

- **`PredictSingle` en bucles**, nunca `Predict` con una sola muestra.
- **Siempre validar `if err != nil`** tras `Fit()`, `Predict()`, `Load()`.
- **Struct embedding:** los modelos como `Lasso` embeben `LinearModel`; los coeficientes se acceden como `model.Coef`.
- **Constructores `New...()`**, nunca structs literal sin inicializar.
- **Cero dependencias nuevas** en el `go.mod` del módulo principal.
- El código generado debe compilar con `go vet ./...` y pasar `go test ./... -race` **antes** de abrir el PR.

---

## 5. Tests

- Todo cambio de comportamiento incluye tests.
- Ejecuta la suite completa antes de enviar: `go test ./... -race`.
- Los tests deben ser deterministas: usa `utils.NewRNG(seed)` en lugar de `rand` global.
- Para cambios de rendimiento, añade un benchmark (`Benchmark*`) en el paquete afectado.

---

## 6. Flujo de pull requests

1. Crea una rama descriptiva: `feature/mi-feature` o `fix/mi-fix`.
2. Un objetivo por PR (mezclar features con fixes complica la revisión).
3. Actualiza la documentación afectada (`README.md`, `USAGE.md`, `docs/`, `CHANGELOG.md`).
4. Comprueba la checklist del [template de PR](.github/PULL_REQUEST_TEMPLATE.md).
5. Describe **qué** cambia, **por qué** cambia y **cómo** se ha validado.

### Convención de commits (recomendada)

```
feat(linear_model): añadir solver alternativo para Ridge
fix(crypto): validar velas con volumen cero
docs(api): actualizar referencia de model_selection
test(ensemble): cubrir PredictProbaSingle
```

---

## 7. Reportar problemas

- **Bugs:** usa la [plantilla de bug report](https://github.com/ateto3/gopherml/issues/new?template=bug_report.yml).
- **Funcionalidades:** usa la [plantilla de feature request](https://github.com/ateto3/gopherml/issues/new?template=feature_request.yml).
- **Vulnerabilidades de seguridad:** no abras una issue pública; sigue [SECURITY.md](SECURITY.md).

---

## 8. Licencia

Al contribuir, otorgas tus contribuciones bajo la [licencia MIT](LICENSE) del proyecto.
