// Copyright 2015 Google Inc. All Rights Reserved.
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

// A tool to measure the upload throughput of GCS.
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/jacobsa/fuse/fsutil"
	"github.com/jacobsa/gcloud/gcs"
	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"
)

var fBucket = flag.String("bucket", "", "Name of bucket.")
var fKeyFile = flag.String("key_file", "", "Path to JSON key file.")
var fSize = flag.Int64("size", 1<<26, "Size of content to write.")

func createBucket() (bucket gcs.Bucket, err error)

func run() (err error) {
	bucket, err := createBucket()
	if err != nil {
		err = fmt.Errorf("createBucket: %v", err)
		return
	}

	// Create a temporary file to hold random contents.
	f, err := fsutil.AnonymousFile("")
	if err != nil {
		err = fmt.Errorf("AnonymousFile: %v", err)
		return
	}

	// Copy a bunch of random data into the file.
	_, err = io.Copy(f, io.LimitReader(rand.Reader, *fSize))
	if err != nil {
		err = fmt.Errorf("Copy: %v", err)
		return
	}

	// Seek back to the start for consumption below.
	_, err = f.Seek(0, 0)
	if err != nil {
		err = fmt.Errorf("Seek: %v", err)
		return
	}

	// Create an object using the contents of the file.
	req := &gcs.CreateObjectRequest{
		Attrs: storage.ObjectAttrs{
			Name: "foo",
		},
		Contents: f,
	}

	_, err = bucket.CreateObject(context.Background(), req)
	if err != nil {
		err = fmt.Errorf("CreateObject: %v", err)
		return
	}

	return
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}
