#!/bin/bash -x

# Copyright 2018 github.com/ucirello
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


openssl genrsa -out fake-ca.key 4096
openssl req -new -x509 \
-subj "/emailAddress=ca@example.com/CN=fake-ca/OU=ops/O=services/L=SomeCity/ST=CA/C=US" \
-days 365 -key fake-ca.key -out fake-ca.crt
openssl genrsa -out fake-server.key 1024
openssl req -new \
-subj "/emailAddress=svr@example.com/CN=fake-svr/OU=ops/O=services/L=SomeCity/ST=CA/C=US" \
-key fake-server.key -out fake-server.csr
openssl x509 -req -days 365 -in fake-server.csr -CA fake-ca.crt -CAkey fake-ca.key \
-extensions v3_server -extfile ./ssl-extensions-x509 \
-set_serial 01 -out fake-server.crt

