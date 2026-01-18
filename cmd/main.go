package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"path/filepath"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/PuerkitoBio/goquery"
)

type Config struct {
	DefaultFile string
	BaseDir     string
	OutDir      string
	DefaultExt  string
}

func loadConfig(path string) Config {
	// Valeurs par défaut
	cfg := Config{
		DefaultFile: "data/input.txt",
		BaseDir:     "data",
		OutDir:     "out",
		DefaultExt: ".txt",
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Config non trouvée, valeurs par défaut utilisées")
		return cfg
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// ignorer commentaires et lignes vides
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "default_file":
			cfg.DefaultFile = value
		case "base_dir":
			cfg.BaseDir = value
		case "out_dir":
			cfg.OutDir = value
		case "default_ext":
			cfg.DefaultExt = value
		}
	}

	return cfg
}

func analyzeFileInfo(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("Erreur : fichier introuvable")
		return
	}

	if info.IsDir() {
		fmt.Println("Erreur : le chemin est un dossier, pas un fichier")
		return
	}

	fmt.Println("\n--- Informations fichier ---")
	fmt.Println("Nom :", info.Name())
	fmt.Println("Taille (octets) :", info.Size())
	fmt.Println("Dernière modification :", info.ModTime())

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Erreur ouverture fichier")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	fmt.Println("Nombre de lignes :", lines)
}

func analyzeWordStats(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Erreur lecture fichier")
		return
	}

	words := strings.Fields(string(data))

	wordCount := 0
	totalLen := 0

	for _, w := range words {
		// ignorer les mots numériques
		if _, err := strconv.Atoi(w); err == nil {
			continue
		}

		wordCount++
		totalLen += len(w)
	}

	fmt.Println("\n--- Statistiques mots ---")
	fmt.Println("Nombre de mots (hors numériques) :", wordCount)

	if wordCount > 0 {
		avg := float64(totalLen) / float64(wordCount)
		fmt.Printf("Longueur moyenne : %.2f\n", avg)
	} else {
		fmt.Println("Longueur moyenne : 0")
	}
}

func filterByKeyword(path string, keyword string, outDir string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Erreur ouverture fichier")
		return
	}
	defer file.Close()

	filteredPath := outDir + "/filtered.txt"
	filteredNotPath := outDir + "/filtered_not.txt"

	filtered, err := os.Create(filteredPath)
	if err != nil {
		fmt.Println("Erreur création filtered.txt")
		return
	}
	defer filtered.Close()

	filteredNot, err := os.Create(filteredNotPath)
	if err != nil {
		fmt.Println("Erreur création filtered_not.txt")
		return
	}
	defer filteredNot.Close()

	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, keyword) {
			fmt.Fprintln(filtered, line)
			count++
		} else {
			fmt.Fprintln(filteredNot, line)
		}
	}

	fmt.Println("\nMot-clé recherché :", keyword)
	fmt.Println("Lignes contenant le mot-clé :", count)
	fmt.Println("Fichiers générés :")
	fmt.Println("-", filteredPath)
	fmt.Println("-", filteredNotPath)
}

func headTail(path string, n int, outDir string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Erreur ouverture fichier")
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		fmt.Println("Fichier vide")
		return
	}

	if n > len(lines) {
		n = len(lines)
	}

	headPath := outDir + "/head.txt"
	tailPath := outDir + "/tail.txt"

	headFile, err := os.Create(headPath)
	if err != nil {
		fmt.Println("Erreur création head.txt")
		return
	}
	defer headFile.Close()

	tailFile, err := os.Create(tailPath)
	if err != nil {
		fmt.Println("Erreur création tail.txt")
		return
	}
	defer tailFile.Close()

	for i := 0; i < n; i++ {
		fmt.Fprintln(headFile, lines[i])
	}

	for i := len(lines) - n; i < len(lines); i++ {
		fmt.Fprintln(tailFile, lines[i])
	}

	fmt.Println("\nHead/Tail générés :")
	fmt.Println("-", headPath)
	fmt.Println("-", tailPath)
}

func listTxtFiles(dir string, ext string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ext) {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files, nil
}

func generateReport(files []string, outDir string) {
	totalFiles := len(files)
	totalLines := 0
	totalWords := 0

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		totalLines += len(lines)

		words := strings.Fields(string(data))
		for _, w := range words {
			if _, err := strconv.Atoi(w); err == nil {
				continue
			}
			totalWords++
		}
	}

	reportPath := outDir + "/report.txt"
	out, err := os.Create(reportPath)
	if err != nil {
		fmt.Println("Erreur création report.txt")
		return
	}
	defer out.Close()

	fmt.Fprintln(out, "=== RAPPORT GLOBAL ===")
	fmt.Fprintln(out, "Fichiers analysés :", totalFiles)
	fmt.Fprintln(out, "Total lignes      :", totalLines)
	fmt.Fprintln(out, "Total mots        :", totalWords)

	fmt.Println("\nRapport généré :", reportPath)
}

func generateIndex(files []string, outDir string) {
	indexPath := outDir + "/index.txt"

	out, err := os.Create(indexPath)
	if err != nil {
		fmt.Println("Erreur création index.txt")
		return
	}
	defer out.Close()

	fmt.Fprintln(out, "=== INDEX DES FICHIERS ===")

	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		fmt.Fprintf(out,
			"Chemin: %s | Taille: %d octets | Modifié: %s\n",
			path,
			info.Size(),
			info.ModTime().Format("2006-01-02 15:04:05"),
		)
	}

	fmt.Println("Index généré :", indexPath)
}

func mergeFiles(files []string, outDir string) {
	mergedPath := outDir + "/merged.txt"

	out, err := os.Create(mergedPath)
	if err != nil {
		fmt.Println("Erreur création merged.txt")
		return
	}
	defer out.Close()

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		
		fmt.Fprintf(out, "\n===== %s =====\n", path)
		out.Write(data)

		
		if len(data) > 0 && data[len(data)-1] != '\n' {
			fmt.Fprintln(out)
		}
	}

	fmt.Println("Fusion générée :", mergedPath)
}

func analyzeWikipedia(article string, outDir string) {
	url := "https://fr.wikipedia.org/wiki/" + article

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Erreur création requête")
		return
	}

	
	req.Header.Set("User-Agent", "FileOps-StudentProject/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur requête HTTP")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Erreur HTTP :", resp.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Erreur parsing HTML")
		return
	}

	var content strings.Builder
	doc.Find("#mw-content-text p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			content.WriteString(text + "\n")
		}
	})

	if content.Len() == 0 {
		fmt.Println("Aucun contenu trouvé")
		return
	}

	
	filename := strings.ReplaceAll(article, "/", "_")
	outPath := outDir + "/wiki_" + filename + ".txt"

	err = os.WriteFile(outPath, []byte(content.String()), 0644)
	if err != nil {
		fmt.Println("Erreur écriture fichier wiki")
		return
	}

	fmt.Println("Article Wikipédia sauvegardé :", outPath)

	
	analyzeWordStats(outPath)
}

func listProcesses(topN int) {
	fmt.Println("\n--- Liste des processus ---")

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		
		cmd = exec.Command("tasklist", "/FO", "CSV")
	} else {
		cmd = exec.Command("ps", "-Ao", "pid,comm")
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Erreur commande système :", err)
		return
	}

	lines := strings.Split(string(out), "\n")

	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fmt.Println(line)
		count++
		if topN > 0 && count >= topN {
			break
		}
	}
}

func searchProcesses(keyword string) {
	fmt.Println("\n--- Recherche processus :", keyword, "---")

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	} else {
		cmd = exec.Command("ps", "-Ao", "pid,comm")
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Erreur commande système :", err)
		return
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
			fmt.Println(line)
		}
	}
}

func killProcess(pid string) {
	fmt.Println("PID à tuer :", pid)
	fmt.Print("Confirmer la terminaison ? (yes/no) : ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	confirm := strings.TrimSpace(scanner.Text())

	if confirm != "yes" {
		fmt.Println("Action annulée.")
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/PID", pid, "/T")
	} else {
		cmd = exec.Command("kill", pid)
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("Erreur kill :", err)
		return
	}

	fmt.Println("Processus terminé.")
}


func main() {
	cfg := loadConfig("config.txt")

	fmt.Println("Configuration chargée :")
	fmt.Println("default_file =", cfg.DefaultFile)
	fmt.Println("base_dir     =", cfg.BaseDir)
	fmt.Println("out_dir      =", cfg.OutDir)
	fmt.Println("default_ext  =", cfg.DefaultExt)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== MENU PRINCIPAL ===")
		fmt.Println("1 - Analyse fichier (Choix A)")
		fmt.Println("2 - Analyse multi-fichiers (Choix B)")
		fmt.Println("3 - Analyser Wikipédia (Choix C)")
		fmt.Println("4 - ProcessOps (Choix D)")
		fmt.Println("0 - Quitter")
		fmt.Print("Votre choix : ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {

		case "1":
	fmt.Println("Analyse du fichier :", cfg.DefaultFile)

	
	analyzeFileInfo(cfg.DefaultFile)

	
	analyzeWordStats(cfg.DefaultFile)

	
	fmt.Print("\nEntrez un mot-clé : ")
	scanner.Scan()
	keyword := strings.TrimSpace(scanner.Text())

	if keyword != "" {
		filterByKeyword(cfg.DefaultFile, keyword, cfg.OutDir)
	}

	
	fmt.Print("\nNombre de lignes pour head/tail : ")
	scanner.Scan()
	nStr := strings.TrimSpace(scanner.Text())

	n, err := strconv.Atoi(nStr)
	if err == nil && n > 0 {
		headTail(cfg.DefaultFile, n, cfg.OutDir)
	} else {
		fmt.Println("Nombre invalide, head/tail ignoré")
	}

	
		case "2":
	fmt.Print("Entrez un dossier (laisser vide pour base_dir) : ")
	scanner.Scan()
	dir := strings.TrimSpace(scanner.Text())

	if dir == "" {
		dir = cfg.BaseDir
	}

	files, err := listTxtFiles(dir, cfg.DefaultExt)
	if err != nil {
		fmt.Println("Erreur lecture dossier :", err)
		break
	}

	fmt.Println("\nFichiers .txt trouvés :", len(files))
	for _, f := range files {
		fmt.Println("-", f)
	}
	generateReport(files, cfg.OutDir)
	generateIndex(files, cfg.OutDir)
	mergeFiles(files, cfg.OutDir)

		case "3":
			fmt.Print("Nom de l'article Wikipédia (ex: Monkey_D._Luffy) : ")
			scanner.Scan()
			article := strings.TrimSpace(scanner.Text())

			if article != "" {
				analyzeWikipedia(article, cfg.OutDir)
			}

		case "4":
	for {
		fmt.Println("\n=== PROCESS OPS ===")
		fmt.Println("1 - Lister les processus (top N)")
		fmt.Println("2 - Rechercher un processus")
		fmt.Println("3 - Tuer un processus")
		fmt.Println("0 - Retour menu principal")
		fmt.Print("Votre choix : ")

		scanner.Scan()
		sub := strings.TrimSpace(scanner.Text())

		switch sub {

		case "1":
			fmt.Print("Top N (0 = tous) : ")
			scanner.Scan()
			nStr := strings.TrimSpace(scanner.Text())
			n, _ := strconv.Atoi(nStr)
			listProcesses(n)

		case "2":
			fmt.Print("Mot-clé : ")
			scanner.Scan()
			kw := strings.TrimSpace(scanner.Text())
			if kw != "" {
				searchProcesses(kw)
			}

		case "3":
			fmt.Print("PID à tuer : ")
			scanner.Scan()
			pid := strings.TrimSpace(scanner.Text())
			if pid != "" {
				killProcess(pid)
			}

		case "0":
			// retour au menu principal
			break

		default:
			fmt.Println("Choix invalide.")
		}

		if sub == "0" {
			break
		}
	}


		case "0":
			fmt.Println("Au revoir.")
			return
		default:
			fmt.Println("Choix invalide, réessayez.")
		}
	}
}

