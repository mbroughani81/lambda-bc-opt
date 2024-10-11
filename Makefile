#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
SHELL:=/bin/bash
OW_USER?=openwhisk
OW_VER?=v1.23:nightly
OW_COMPILER?=$(OW_USER)/action-golang-$(OW_VER)

HANDLERS_DIR=lambda-bc-opt/handlers/wsk
ZIP_DIR=zipzip

.PHONY: deploy devel clean $(ZIP_DIR)

# Create an action for each subfolder in $(HANDLERS_DIR)
devel:
	rm -rf $(ZIP_DIR); \
	for dir in $(HANDLERS_DIR)/*/; do \
		zip_name=$$(basename $$dir).zip; \
		action_name=$$(basename $$dir); \
		echo $$zip_name; \
		rm -rf $(ZIP_DIR); \
		mkdir -p $(ZIP_DIR); \
		cp $$dir/main.go $(ZIP_DIR); \
		cp -r lambda-bc-opt/* $(ZIP_DIR); \
		(cd $(ZIP_DIR) && zip -r $$zip_name . && \
		wsk action update $$action_name $$zip_name --main main --docker $(OW_COMPILER) --memory 512; \
		); \
	done; \
	rm -rf $(ZIP_DIR);

# Deploy to OpenWhisk
deploy:
	rm -rf $(ZIP_DIR)
	for dir in $(HANDLERS_DIR)/*/; do \
		rm -rf $(ZIP_DIR); \
		action_name=$$(basename $$dir); \
		echo "Building and zipping $$action_name"; \
		mkdir -p $(ZIP_DIR); \
		cp $$dir/main.go $(ZIP_DIR); \
		cp -r lambda-bc-opt/* $(ZIP_DIR); \
		cd $(ZIP_DIR) && zip - -r . | docker run -i $(OW_COMPILER) -compile main > ../$$action_name.zip && cd .. ;\
		echo "Deploying $$action_name"; \
		wsk action update $$action_name $$action_name.zip --main exec --docker $(OW_COMPILER) --memory 128; \
		rm $$action_name.zip; \
	done;

# Copy the main.go file from the subfolder and zip it
zip:
	rm -rf $(ZIP_DIR)
	mkdir -p $(ZIP_DIR)
	cp $(DIR)/main.go $(ZIP_DIR)
	(cd $(ZIP_DIR) && zip -r ../$(SRCZIP) .)

clean:
	rm -rf $(ZIP_DIR)
