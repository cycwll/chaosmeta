/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package namespace

import (
	"chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

type Namespace struct {
	Id          int    `json:"id" orm:"pk;auto;column(id)"`
	Name        string `json:"name" orm:"column(name); size(255);index"`
	Description string `json:"description" orm:"column(description); size(1024)"`
	Creator     int    `json:"creator" orm:"column(creator); index"`
	IsDefault   bool   `json:"is_default" orm:"column(is_default)"`
	models.BaseTimeModel
}

func (u *Namespace) TableName() string {
	return "namespace"
}

func InsertNamespace(ctx context.Context, namespace *Namespace) (int64, error) {
	if namespace == nil {
		return 0, errors.New("namespace is nil")
	}
	id, err := models.GetORM().Insert(namespace)
	return id, err
}

func UpdateNamespace(ctx context.Context, namespace *Namespace) (int64, error) {
	if namespace == nil {
		return 0, errors.New("namespace is nil")
	}
	num, err := models.GetORM().Update(namespace)
	return num, err
}

func DeleteNamespace(ctx context.Context, id int) (int64, error) {
	num, err := models.GetORM().Delete(&Namespace{Id: id})
	return num, err
}

func GetNamespaceByName(ctx context.Context, namespace *Namespace) error {
	if namespace == nil {
		return errors.New("namespace is nil")
	}
	return models.GetORM().Read(namespace, "name")
}

func GetNamespaceById(ctx context.Context, namespace *Namespace) error {
	if namespace == nil {
		return errors.New("namespace is nil")
	}
	return models.GetORM().Read(namespace)
}

func GetDefaultNamespace(ctx context.Context, namespace *Namespace) error {
	if namespace == nil {
		return errors.New("namespace is nil")
	}
	return models.GetORM().QueryTable(namespace.TableName()).Filter("is_default", true).One(namespace)
}

func GetAllNamespaces() ([]*Namespace, error) {
	var namespaces []*Namespace
	_, err := models.GetORM().QueryTable("namespace").All(&namespaces)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}

func QueryNamespaces(ctx context.Context, name, creator, orderBy string, page, pageSize int) (int64, []Namespace, error) {
	ns, namespaceList := Namespace{}, new([]Namespace)
	querySeter := models.GetORM().QueryTable(ns.TableName())
	namespaceQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	if len(name) > 0 {
		namespaceQuery.Filter("name", models.CONTAINS, true, name)
	}

	if len(creator) > 0 {
		namespaceQuery.Filter("creator", models.NEGLECT, false, creator)
	}

	orderByList := []string{}
	if orderBy != "" {
		orderByList = append(orderByList, orderBy)
	} else {
		orderByList = append(orderByList, "id")
	}
	namespaceQuery.OrderBy(orderByList...)

	if err := namespaceQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	totalCount, err := namespaceQuery.GetOamQuerySeter().All(namespaceList)
	if err == orm.ErrNoRows {
		return 0, nil, nil
	}
	return totalCount, *namespaceList, err
}

func ListNamespaces(ctx context.Context, namespaceId []int, name string, creator, orderBy string, page, pageSize int) (int64, []Namespace, error) {
	ns, namespaceList := Namespace{}, new([]Namespace)
	querySeter := models.GetORM().QueryTable(ns.TableName())
	namespaceQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	if len(namespaceId) > 0 {
		namespaceQuery.Filter("id", models.IN, false, namespaceId)
	}
	if len(name) > 0 {
		namespaceQuery.Filter("name", models.CONTAINS, true, name)
	}

	if len(creator) > 0 {
		namespaceQuery.Filter("creator", models.NEGLECT, false, creator)
	}

	orderByList := []string{}
	if orderBy != "" {
		orderByList = append(orderByList, orderBy)
	} else {
		orderByList = append(orderByList, "id")
	}
	namespaceQuery.OrderBy(orderByList...)

	if err := namespaceQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}
	totalCount, _ := namespaceQuery.GetOamQuerySeter().Count()

	_, err = namespaceQuery.GetOamQuerySeter().All(namespaceList)
	if err == orm.ErrNoRows {
		return 0, nil, nil
	}
	return totalCount, *namespaceList, err
}
