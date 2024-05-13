package services

import (
	"bytes"
	"github.com/SiriusServiceDesk/notification-service/internal/models"
	"github.com/SiriusServiceDesk/notification-service/internal/repository"
	tmpl "html/template"
)

type TemplateService interface {
	Create(template *models.Template) error
	GetAll() ([]*models.Template, error)
	GetByName(name string) (*models.Template, error)
	Update(template *models.Template) error
	Delete(id string) error
	Render(templateName string, messageData models.JSONB) (string, error)
}

type TemplateServiceImpl struct {
	repos    repository.TemplateRepository
	template *models.Template
}

func (t *TemplateServiceImpl) Create(template *models.Template) error {
	return t.repos.CreateTemplate(template)
}

func (t *TemplateServiceImpl) GetAll() ([]*models.Template, error) {
	return t.repos.GetTemplates()
}

func (t *TemplateServiceImpl) GetByName(name string) (*models.Template, error) {
	return t.repos.GetTemplate(name)
}

func (t *TemplateServiceImpl) Update(template *models.Template) error {
	return t.repos.UpdateTemplate(template)
}

func (t *TemplateServiceImpl) Delete(id string) error {
	return t.repos.DeleteTemplate(id)
}

func (t *TemplateServiceImpl) Render(templateName string, messageData models.JSONB) (string, error) {
	template, err := t.repos.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	tmp, err := tmpl.New("emailTemplate").Parse(template.Html)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err = tmp.Execute(&tplBuffer, messageData); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}

func NewTemplateService(repos repository.TemplateRepository) *TemplateServiceImpl {
	return &TemplateServiceImpl{repos: repos}
}
