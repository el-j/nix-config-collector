package mapper

import (
	"strings"

	"github.com/el-j/nix-config-collector/internal/core/models"
)

// BrewToNixMapping maps common Homebrew formula names to nixpkgs attribute names
var BrewToNixMapping = map[string]string{
	"git":           "git",
	"curl":          "curl",
	"wget":          "wget",
	"ripgrep":       "ripgrep",
	"fd":            "fd",
	"fzf":           "fzf",
	"jq":            "jq",
	"htop":          "htop",
	"tmux":          "tmux",
	"neovim":        "neovim",
	"vim":           "vim",
	"bat":           "bat",
	"exa":           "eza",
	"eza":           "eza",
	"lsd":           "lsd",
	"starship":      "starship",
	"fish":          "fish",
	"zsh":           "zsh",
	"bash":          "bash",
	"python3":       "python3",
	"python@3.11":   "python311",
	"python@3.12":   "python312",
	"node":          "nodejs",
	"go":            "go",
	"rust":          "rustup",
	"cmake":         "cmake",
	"make":          "gnumake",
	"gcc":           "gcc",
	"llvm":          "llvm",
	"openssl":       "openssl",
	"readline":      "readline",
	"sqlite":        "sqlite",
	"postgresql":    "postgresql",
	"mysql":         "mysql",
	"redis":         "redis",
	"nginx":         "nginx",
	"docker":        "docker",
	"kubectl":       "kubectl",
	"helm":          "kubernetes-helm",
	"terraform":     "terraform",
	"awscli":        "awscli2",
	"gh":            "gh",
	"hub":           "hub",
	"git-lfs":       "git-lfs",
	"imagemagick":   "imagemagick",
	"ffmpeg":        "ffmpeg",
	"yt-dlp":        "yt-dlp",
	"aria2":         "aria2",
	"tree":          "tree",
	"ncdu":          "ncdu",
	"tldr":          "tealdeer",
	"mkcert":        "mkcert",
	"direnv":        "direnv",
	"nix":           "nix",
	"age":           "age",
	"gnupg":         "gnupg",
	"pass":          "pass",
	"1password-cli": "_1password-cli",
	"mas":           "mas",
	"stow":          "stow",
	"zoxide":        "zoxide",
	"atuin":         "atuin",
	"delta":         "delta",
	"difftastic":    "difftastic",
	"lazygit":       "lazygit",
	"lazydocker":    "lazydocker",
	"k9s":           "k9s",
	"yq":            "yq-go",
	"xq":            "xq",
	"duf":           "duf",
	"glow":          "glow",
	"charm":         "charm",
	"carapace":      "carapace",
	"nushell":       "nushell",
	"helix":         "helix",
	"zellij":        "zellij",
	"wezterm":       "wezterm",
	"alacritty":     "alacritty",
	"kitty":         "kitty",
	"exercism":      "exercism",
	"protobuf":      "protobuf",
	"grpc":          "grpc",
	"buf":           "buf",
}

// Mapper handles mapping between macOS packages and Nix equivalents
type Mapper struct{}

// New creates a new Mapper
func New() *Mapper {
	return &Mapper{}
}

// MapPackages maps packages to their Nix equivalents
func (m *Mapper) MapPackages(packages []models.Package) ([]models.Package, error) {
	mapped := make([]models.Package, len(packages))
	for i, pkg := range packages {
		mapped[i] = pkg
		if pkg.Type == models.PackageTypeBrewFormula {
			nixName, ok := BrewToNixMapping[strings.ToLower(pkg.Name)]
			if ok {
				mapped[i].NixName = nixName
			}
		}
	}
	return mapped, nil
}

// MapServices maps services (no-op for now, returns as-is)
func (m *Mapper) MapServices(services []models.Service) ([]models.Service, error) {
	return services, nil
}

// LookupNixName returns the Nix package name for a brew formula
func LookupNixName(brewName string) (string, bool) {
	name, ok := BrewToNixMapping[strings.ToLower(brewName)]
	return name, ok
}

// GetUnmappedPackages returns packages without a Nix equivalent
func GetUnmappedPackages(packages []models.Package) []models.Package {
	var unmapped []models.Package
	for _, pkg := range packages {
		if pkg.Type == models.PackageTypeBrewFormula && pkg.NixName == "" {
			unmapped = append(unmapped, pkg)
		}
	}
	return unmapped
}

// SplitByDestination splits packages into nix packages vs homebrew packages
func SplitByDestination(packages []models.Package) (nixPkgs, brewFormulae, brewCasks []models.Package, masApps []models.Package) {
	for _, pkg := range packages {
		switch pkg.Type {
		case models.PackageTypeBrewFormula:
			if pkg.NixName != "" {
				nixPkgs = append(nixPkgs, pkg)
			} else {
				brewFormulae = append(brewFormulae, pkg)
			}
		case models.PackageTypeBrewCask:
			brewCasks = append(brewCasks, pkg)
		case models.PackageTypeMAS:
			masApps = append(masApps, pkg)
		}
	}
	return
}
