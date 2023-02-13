/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import "github.com/a-clap/iot/internal/embedded/models"

type GlobalConfig struct {
	heaters []models.HeaterConfig
}

type PhaseConfig struct {
	heaters []models.HeaterConfig
}
