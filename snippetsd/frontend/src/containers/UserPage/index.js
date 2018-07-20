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
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import {
  Button,
  Col,
  Form,
  Grid,
  PageHeader,
  Row,
  FormControl
} from 'react-bootstrap'
import { loadSnippetsByUser, saveSnippet } from './actions'
import groupBy from 'lodash/groupBy'
import moment from 'moment'

import './style.css'

class SubmitSnippetPage extends React.Component {
  constructor (props) {
    super(props)

    this.setContent = this.setContent.bind(this)
    this.submit = this.submit.bind(this)

    this.state = {
      content: ''
    }
  }

  componentDidMount () {
    this.props.loadSnippetsByUser()
  }

  setContent (e) {
    e.preventDefault()
    var content = e.target.value
    this.setState({ content })
  }
  submit (e) {
    e.preventDefault()
    this.props.saveSnippet(this.state.content)
  }

  render () {
    var snippets = groupBy(this.props.snippets, function (v) {
      return v.week_start
    })

    return (
      <Grid className='user-snippet-grid'>
        <Row>
          <Col>
            <div className='user-snippet-container'>
              <PageHeader> What did you do past week? </PageHeader>

              <Form onSubmit={this.submit}>
                <FormControl componentClass='textarea' className='user-snippet-content' onChange={this.setContent} />
                <div className='user-snippet-submit'><Button type='submit'>submit</Button></div>
              </Form>
            </div>
          </Col>
        </Row>
        <Row>
          <Col>
            <PageHeader> Past Snippets </PageHeader>
            {
              Object.entries(snippets).map((week) => (
                <div key={week[0]} className='user-past-snippet'>
                  <strong>Week starting {moment(week[0]).format('MMMM Do YYYY')}: </strong>
                  { week[1].map((snippet) => (<div key={snippet.user.email}>{snippet.contents}</div>)) }
                </div>
              ))
            }
          </Col>
        </Row>
      </Grid>
    )
  }
}

const s2p = state => ({ snippets: state.snippets.snippets })
const d2p = dispatch => bindActionCreators({
  loadSnippetsByUser,
  saveSnippet
}, dispatch)
export default connect(s2p, d2p)(SubmitSnippetPage)
