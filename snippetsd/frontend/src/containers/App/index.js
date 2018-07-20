// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from 'react'
import { Route, Link } from 'react-router-dom'
import { Nav, NavItem, Navbar } from 'react-bootstrap'
import WeekPage from '../WeekPage'
import UserPage from '../UserPage'

const App = () => (
  <div>
    <Navbar>
      <Navbar.Header>
        <Navbar.Brand>
          <Link to='/view'>Weekly snippets</Link>
        </Navbar.Brand>
      </Navbar.Header>
      <Nav>
        <NavItem componentClass={Link} eventKey={1}
          href='/' to='/'>my snippets</NavItem>
      </Nav>
    </Navbar>

    <main>
      <Route exact path='/view' component={WeekPage} />
      <Route exact path='/' component={UserPage} />
    </main>
  </div>
)

export default App
