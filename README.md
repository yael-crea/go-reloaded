# go-reloaded

Programme Go en ligne de commande qui lit un fichier texte, applique les règles de correction du sujet, puis écrit le résultat dans un autre fichier.

## Exécuter

```powershell
go run . sample.txt result.txt
```

## Tester

```powershell
go test ./... -v
go vet ./...
gofmt -l .
```
