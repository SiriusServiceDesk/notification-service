package repository

import (
	"errors"
	"fmt"
	"github.com/SiriusServiceDesk/notification-service/internal/config"
	"github.com/SiriusServiceDesk/notification-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TemplateRepository interface {
	GetTemplates() ([]*models.Template, error)
	GetTemplate(name string) (*models.Template, error)
	CreateTemplate(template *models.Template) error
	UpdateTemplate(template *models.Template) error
	DeleteTemplate(id string) error
}

type TemplateRepositoryImpl struct {
	db *gorm.DB
}

func (r *TemplateRepositoryImpl) GetTemplates() ([]*models.Template, error) {
	var templates []*models.Template
	result := r.db.Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}
	return templates, nil
}

func (r *TemplateRepositoryImpl) GetTemplate(name string) (*models.Template, error) {
	var template *models.Template
	result := r.db.Where("template_name =?", name).First(&template)
	if result.Error != nil {
		return nil, result.Error
	}
	return template, nil
}

func (r *TemplateRepositoryImpl) CreateTemplate(template *models.Template) error {
	result := r.db.Create(&template)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TemplateRepositoryImpl) UpdateTemplate(template *models.Template) error {
	result := r.db.Save(&template)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TemplateRepositoryImpl) DeleteTemplate(id string) error {
	result := r.db.Delete(&models.Template{}, "id =?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TemplateRepositoryImpl) seeds() {
	const serverStartName = "serverStartNotification"
	const successfulRegistrationName = "successfulRegistration"
	const verifyingEmailName = "verifyingEmail"
	const resetPassword = "resetPassword"
	const updateApplications = "updateApplication"

	templates := []models.Template{
		{
			TemplateName: serverStartName,
			Html:         `<p> Сервис {{.ServiceName}}</p> <br /> Запущен успешно!`,
		},
		{
			TemplateName: successfulRegistrationName,
			Html:         `Congratulations! Welcome to our website, {{.Name}} {{.Surname}}`,
		},
		{
			TemplateName: verifyingEmailName,
			Html: `<div class="container" style="display:flex;justify-content:center; align-items: center;flex-direction: column;gap:15px;width:100%; font-family: sans-serif;padding: 14px 20px; font-size:13px;font-weight500">
						  <div style="margin-left: auto">сириусдеск.рф</div>
						  <div style="width:100%;display:flex; justify-content:center; padding: 20px 15px;background-color:#5046E6;border-radius:10px 10px 10px 0px;color:#FFF">Спасибо за регистрацию на сириусдеск.рф!</div>
						  <div style="width:100%; padding: 20px 15px;background-color:#F3F3F3;border-radius:10px 10px 10px 0px;">
								Ваш код для подтверждения email:&nbsp;<b style="color:#5046E6;">{{.Code}}</b></div>
						  </div>`,
		},
		{
			TemplateName: resetPassword,
			Html:         `Тут ссылка для сброса пароля по идее пока не делал`,
		},
		{
			TemplateName: updateApplications,
			Html: `<div class="container" style="display:flex;justify-content:center; align-items: center;flex-direction: column;gap:15px;width:100%; font-family: sans-serif;padding: 14px 20px; font-size:12px;font-weight500">	<div style="margin-left: auto">сириусдеск.рф</div>
  					<div style="width:100%; padding: 20px 15px;background-color:#5046E6;border-radius:10px 10px 10px 0px;color:#FFF">Вашей заявке {{.AppId}} был присвоен статус <b>{{.Status}}</b></div>
  					{{if .Comment}}<div style="width:100%; padding: 20px 15px;background-color:#F3F3F3;border-radius:10px 10px 10px 0px;"><b>Исполнитель оставил комментарий</b> “{{.Comment}}”</div>
					{{else}} <div style="width:100%; padding: 20px 15px;background-color:#F3F3F3;border-radius:10px 10px 10px 0px;"><b>Исполнитель не оставил комментарий :( </b></div>
					{{end}}</div>
			`,
		},
	}

	for _, template := range templates {
		_, err := r.GetTemplate(template.TemplateName)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.CreateTemplate(&template)
		}
	}

}

func NewTemplateRepository() TemplateRepository {
	cfg := config.GetConfig().Db
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	pgSvc := &TemplateRepositoryImpl{db: db}
	err = db.AutoMigrate(&models.Template{})
	if err != nil {
		panic(err)
	}

	pgSvc.seeds()

	return pgSvc
}
