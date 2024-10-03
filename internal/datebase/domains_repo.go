package datebase

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sunwild/api/pkg/domains"
)

func (r *Repository) GetAllDomains() ([]*domains.Domain, error) {
	domainsList := make([]*domains.Domain, 0)
	rows, err := r.db.Query("SELECT * FROM domains")

	defer rows.Close()

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения списка доменов")
	}

	for rows.Next() {
		var domainsIter domains.Domain
		// Считываем данные в структуру domain
		err := rows.Scan(&domainsIter.ID, &domainsIter.Name)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка при сканировании данных домена")
		}
		// Добавляем структуру в список
		domainsList = append(domainsList, &domainsIter)
	}

	return domainsList, nil
}

func (r *Repository) GetDomainById(domainId int) (*domains.Domain, error) {
	var domain domains.Domain
	err := r.db.QueryRow("SELECT id, name FROM domains WHERE id = ?", domainId).Scan(&domain.ID, &domain.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("domain not found") // Возвращаем ошибку
		}
		return nil, err // Возвращаем ошибку, если произошла другая проблема
	}
	return &domain, nil
}

func (r *Repository) AddDomain(domain *domains.Domain) error {
	_, err := r.db.Exec("INSERT INTO domains(name) VALUES(?)", domain.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteDomainById(domainId int) error {
	_, err := r.db.Exec("DELETE FROM domains WHERE id = ?", domainId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateDomainById(domain *domains.Domain) error {
	_, err := r.db.Exec("UPDATE domains SET name = ? WHERE id = ?", domain.Name, domain.ID)
	if err != nil {
		return err
	}
	return nil
}
