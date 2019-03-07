// Copyright 2018 github.com/ucirello
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

import config from '../../config'

var cfg = config()

export function initialDataload () {
  return (dispatch) => {
    fetch(cfg.http + '/state', {
      credentials: 'same-origin'
    })
    .then(res => res.json())
    .catch((e) => {
      console.log('cannot load link state:', e)
    })
    .then((bookmarks) => {
      dispatch({
        type: 'INITIAL_LOAD',
        bookmarks: bookmarks
      })
    })
  }
}

export function deleteBookmark (id) {
  return (dispatch) => {
    fetch(cfg.http + '/deleteBookmark', {
      method: 'POST',
      body: JSON.stringify({ id }),
      credentials: 'same-origin'
    })
    .then(res => res.json())
    .catch((e) => {
      console.log('cannot delete bookmark:', e)
    })
  }
}

export function markBookmarkAsRead (id) {
  return (dispatch) => {
    fetch(cfg.http + '/markBookmarkAsRead', {
      method: 'POST',
      body: JSON.stringify({ id }),
      credentials: 'same-origin'
    })
    .then(res => res.json())
    .catch((e) => {
      console.log('cannot mark bookmark as read:', e)
    })
  }
}
