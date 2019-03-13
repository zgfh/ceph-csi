/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/klog"
)

// CreatePersistanceStorage creates storage path and initializes new cache
func CreatePersistanceStorage(sPath, metaDataStore, driverName string) (CachePersister, error) {
	var err error
	if err = createPersistentStorage(path.Join(sPath, "controller")); err != nil {
		klog.Errorf("failed to create persistent storage for controller: %v", err)
		return nil, err
	}

	if err = createPersistentStorage(path.Join(sPath, "node")); err != nil {
		klog.Errorf("failed to create persistent storage for node: %v", err)
		return nil, err
	}

	cp, err := NewCachePersister(metaDataStore, driverName)
	if err != nil {
		klog.Errorf("failed to define cache persistence method: %v", err)
		return nil, err
	}
	return cp, err
}

func createPersistentStorage(persistentStoragePath string) error {
	return os.MkdirAll(persistentStoragePath, os.FileMode(0755))
}

// ValidateDriverName validates the driver name
func ValidateDriverName(driverName string) error {
	if len(driverName) == 0 {
		return errors.New("driver name is empty")
	}

	if len(driverName) > 63 {
		return errors.New("driver name length should be less than 63 chars")
	}
	var err error
	for _, msg := range validation.IsDNS1123Subdomain(strings.ToLower(driverName)) {
		if err == nil {
			err = errors.New(msg)
			continue
		}
		err = errors.Wrap(err, msg)
	}
	return err
}
