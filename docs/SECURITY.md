# 🔒 Seguridad

## Versiones soportadas

| Versión | Soportada |
| :--- | :--- |
| Última versión en `main` | ✅ |
| Cualquier versión anterior | ❌ |

Este proyecto se distribuye sin garantías (ver [LICENSE](LICENSE)). Solo se corrigen vulnerabilidades en la rama principal.

## Reportar una vulnerabilidad

Por favor, **no abras una issue pública** para reportar vulnerabilidades de seguridad.

Usa la función de [Security Advisories de GitHub](https://github.com/ateto3/gopherml/security/advisories/new) para informar de forma privada. Si no puedes usarla, contacta directamente al mantenedor a través de GitHub ([@ateto3](https://github.com/ateto3)).

Incluya en el reporte:

1. **Descripción** del problema y su impacto potencial.
2. **Pasos para reproducirlo** (pseudocódigo, snippet o test mínimo).
3. **Versión / commit** afectado.
4. **Mitigación conocida**, si la hay.

## Qué esperar

- **Acuse de recibo** en un plazo razonable (objetivo: 72 horas).
- **Evaluación** del impacto y confirmación de la vulnerabilidad.
- **Corrección y publicación** antes de cualquier divulgación coordinada.
- Crédito en el `CHANGELOG.md` si lo deseas (o anonimato si lo prefieres).

## Alcance

Dentro del alcance:

- Validación de datos (`utils.CheckFinite`) y errores de seguridad en la gestión de entradas.
- Persistencia (`persistence.Load`): corrupción de archivos, validación de integridad de modelos.
- Cualquier fuga de memoria o condición de carrera detectable con `go test -race`.

Fuera de alcance:

- Problemas de rendimiento sin impacto de seguridad.
- Vulnerabilidades en exchanges o APIs de terceros usadas junto a GopherML.
- `go-nlp/nltk/`: código de referencia importado de NLTK (Python) para comparación algorítmica; no forma parte del runtime de Go.
