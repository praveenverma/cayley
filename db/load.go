// Copyright 2014 The Cayley Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"io"

	"github.com/google/cayley/config"
	"github.com/google/cayley/graph"
	"github.com/google/cayley/quad"
)

func Load(ts graph.TripleStore, cfg *config.Config, dec quad.Unmarshaler) error {
	bulker, canBulk := ts.(graph.BulkLoader)
	if canBulk {
		switch err := bulker.BulkLoad(dec); err {
		case nil:
			return nil
		case graph.ErrCannotBulkLoad:
			// Try individual loading.
		default:
			return err
		}
	}

	block := make([]quad.Quad, 0, cfg.LoadSize)
	for {
		t, err := dec.Unmarshal()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		block = append(block, t)
		if len(block) == cap(block) {
			ts.AddTripleSet(block)
			block = block[:0]
		}
	}
	ts.AddTripleSet(block)

	return nil
}
