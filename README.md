# BOLLON Hugo

Versionné et hébergé ici: [https://github.com/hbollon/go-hash](https://github.com/hbollon/go-hash)

Nécessite [Go 1.17](https://golang.org/doc/install) pour être exécuté avec la commande `go run ./cmd`

Une **TUI** a été réalisé à l'aide de la librairie [Bubbletea](https://github.com/charmbracelet/bubbletea), les commandes pour la navigation au sein de l'interface sont indiquées constament en bas de cette dernière.
Cette dernière devrait être compatible dans n'importe quel terminal et à en tout cas été testé:

- Sous Linux (Manjaro) avec le combo **Konsole** et **zsh**
- Sous Windows avec **cmd** et **Windows Terminal**

Pour pouvoir calculer le coverage et tenter un crack de hash il faut d'abord charger une table en mémoire. Pour ce faire, il vous suffit d'en générer une nouvelle ou d'en charger une existante via la TUI.

La fonction invert semble ne pas fonctionner parfaitement car avec un grand alphabet (environ supérieur à 16 caractères) il a du mal à trouver certains mots alors que le coverage est au dessus de 90%.
Je pense que cela proviens des valeurs aléatoires insérés lors de la génération de la table.

## TP1 : go-hash

### Question 7

Le paramètre _*t*_ permet d'éviter les collisions et donc d'augmenter la couverture de la table.
