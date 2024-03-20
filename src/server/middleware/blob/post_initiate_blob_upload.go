// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blob

import (
	"net/http"
	"time"

	"github.com/goharbor/harbor/src/lib/log"
	"github.com/goharbor/harbor/src/pkg/distribution"
	"github.com/goharbor/harbor/src/server/middleware"
)

// PostInitiateBlobUploadMiddleware middleware to add blob to project after mount blob success
func PostInitiateBlobUploadMiddleware() func(http.Handler) http.Handler {
	return middleware.AfterResponse(func(w http.ResponseWriter, r *http.Request, statusCode int) error {
		if statusCode != http.StatusCreated {
			return nil
		}
		now := time.Now()
		ctx := r.Context()
		logger := log.G(ctx).WithFields(log.Fields{"middleware": "blob", "action": "post initiate blob_upload", "url": r.URL.Path})
		defer func() {
			logger.Infof("**post initiate blob upload with method %s take: %s", r.Method, time.Since(now))
		}()

		query := r.URL.Query()

		mount := query.Get("mount")
		if mount == "" {
			return nil
		}

		project, err := projectController.GetByName(ctx, distribution.ParseProjectName(r.URL.Path))
		if err != nil {
			logger.Errorf("get project failed, error: %v", err)
			return err
		}

		if err := blobController.AssociateWithProjectByDigest(ctx, mount, project.ProjectID); err != nil {
			logger.Errorf("mount blob %s to project %s failed, error: %v", mount, project.Name, err)
			return err
		}

		return nil
	})
}
