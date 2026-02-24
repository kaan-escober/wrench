package ui

import (
	"fmt"
	"strings"

	"github.com/kaan-escober/wrench/internal/config"
	"github.com/kaan-escober/wrench/internal/theme"
)

func (m Model) viewBYOK() string {
	switch m.byokStep {

	case WizProvider:
		return viewHeader("PROVIDER", "Choose a provider or select an existing group") +
			m.providerList.render(true)

	case WizGroupDetail:
		g := m.providerGroups[m.currentGroupIdx]
		name := groupDisplayName(g)
		return viewHeader("PROVIDER GROUP", name+" · "+fmt.Sprintf("%d model(s)", len(g.Models))) +
			m.detailList.render(true)

	case WizModelEdit:
		dn := m.editingModel.DisplayName
		if dn == "" {
			dn = m.editingModel.Model
		}
		return viewHeader("EDIT MODEL", dn) +
			m.detailList.render(true)

	case WizModelField:
		return m.viewModelField()

	case WizURL:
		return viewHeader("BASE URL", "e.g. https://openrouter.ai/api/v1") +
			theme.PromptStr() + m.textInput.View()

	case WizTitle:
		return viewHeader("DISPLAY NAME", "Shown next to your models in Droid's model selector") +
			theme.PromptStr() + m.textInput.View() + "\n" +
			theme.Muted.Render("  Extracted: "+m.extractedName)

	case WizKey:
		return viewHeader("API KEY", "Supports raw keys and ${ENV_VAR} references") +
			theme.PromptStr() + m.textInput.View() + "\n" +
			theme.Muted.Render("  Stored in: "+config.SettingsPath())

	case WizFetching:
		return viewHeader("FETCHING", "") +
			theme.Muted.Render(m.spinner.View()+" Fetching models from ") +
			theme.Accent.Render(m.displayTitle) + "..."

	case WizModels:
		if len(m.availableModels) == 0 {
			return viewHeader("MODEL ID", "Could not auto-fetch — enter a model ID manually") +
				theme.Muted.Render("  e.g. gpt-4o, claude-opus-4-5, qwen3:4b") + "\n\n" +
				theme.PromptStr() + m.textInput.View()
		}
		count := len(m.modelList.selectedValues())
		sel := ""
		if count > 0 {
			sel = "  " + theme.BadgeSuccess.Render(fmt.Sprintf(" %d selected ", count))
		}
		return viewHeader("SELECT MODELS", "space · toggle   enter · confirm"+sel) +
			m.modelList.render(true)

	case WizSettingsTokens:
		return viewHeader("MAX TOKENS", "Default: 16384  ·  Most providers: 8192–131072") +
			theme.PromptStr() + m.textInput.View()

	case WizSettingsImages:
		return viewHeader("IMAGE SUPPORT", "Can these models process image inputs?") +
			m.detailList.render(true)

	case WizConfirm:
		return viewHeader("CONFIRM", "Review before saving to ~/.factory/settings.json") +
			m.renderBYOKSummary() + "\n\n" +
			m.detailList.render(true)

	case WizSaving:
		return viewHeader("SAVING", "") +
			theme.Muted.Render(m.spinner.View()+" Writing configuration...")

	case WizDone:
		lines := make([]string, len(m.selectedModels))
		for i, id := range m.selectedModels {
			dn := m.modelDisplayNames[id]
			if dn == "" {
				dn = id
			}
			lines[i] = theme.Success.Render("  ● ") + theme.Primary.Render(dn)
		}
		return viewHeader("DONE", "") +
			theme.BadgeSuccess.Render(" SAVED ") + "\n\n" +
			theme.Primary.Render(fmt.Sprintf("%d model(s) added to Factory", len(m.selectedModels))) + "\n" +
			strings.Join(lines, "\n") + "\n" +
			theme.Muted.Render("  → "+config.SettingsPath()) + "\n\n" +
			m.detailList.render(true)
	}

	return ""
}

func (m Model) viewModelField() string {
	switch m.editFieldKey {
	case "displayName":
		return viewHeader("DISPLAY NAME", "") + theme.PromptStr() + m.textInput.View()
	case "model":
		return viewHeader("MODEL ID", "") + theme.PromptStr() + m.textInput.View()
	case "baseUrl":
		return viewHeader("BASE URL", "") + theme.PromptStr() + m.textInput.View()
	case "apiKey":
		return viewHeader("API KEY", "") + theme.PromptStr() + m.textInput.View()
	case "maxOutputTokens":
		return viewHeader("MAX TOKENS", "") + theme.PromptStr() + m.textInput.View()
	case "provider":
		return viewHeader("API TYPE", "Select the API protocol") + m.detailList.render(true)
	case "supportsImages":
		return viewHeader("IMAGE SUPPORT", "") + m.detailList.render(true)
	case "delete":
		return viewHeader("DELETE MODEL", "Are you sure?") + m.detailList.render(true)
	}
	return ""
}

func (m Model) renderBYOKSummary() string {
	row := func(label, value string) string {
		l := theme.Muted.Render(fmt.Sprintf("  %-14s", label+":"))
		return l + theme.Primary.Render(value)
	}

	models := m.selectedModels
	modelNames := make([]string, len(models))
	for i, id := range models {
		dn := m.modelDisplayNames[id]
		if dn == "" {
			dn = id
		}
		modelNames[i] = dn
	}

	images := "No"
	if m.supportsImages {
		images = "Yes"
	}

	return strings.Join([]string{
		row("Provider", m.displayTitle),
		row("Base URL", m.baseURL),
		row("Type", m.providerType),
		row("Models", strings.Join(modelNames, ", ")),
		row("Max tokens", fmt.Sprintf("%d", m.maxOutputTokens)),
		row("Images", images),
	}, "\n")
}
