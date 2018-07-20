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
import { loadSnippets } from './actions'
import { Col, Grid, PageHeader, Panel, Row } from 'react-bootstrap'
import groupBy from 'lodash/groupBy'
import moment from 'moment'

import './style.css'

export class WeekPage extends React.PureComponent {
  componentDidMount () {
    this.props.loadSnippets()
  }
  render () {
    var snippets = groupBy(this.props.snippets, function (v) {
      return v.week_start
    })

    let groupedSnippets = {}
    for (let i in snippets) {
      groupedSnippets[i] = groupBy(snippets[i], function (v) {
        return v.user.team
      })
    }

    console.log(groupedSnippets)

    return (
      <Grid>
        {
            Object.entries(groupedSnippets).map((week) => (
              <Row key={week[0]}>
                <Col>
                  <PageHeader> Snippets for the week starting {moment(week[0]).format('MMMM Do YYYY')}: </PageHeader>
                  {
                    Object.entries(week[1]).map((team) => (
                      <Panel key={team[0]} className='snippet-card'>
                        <strong>{team[0]}</strong>
                        {team[1].map(
                          (snippet) => (
                            <div key={snippet.user.email} className='user-snippet'>
                              <span>{snippet.user.email}</span>
                              <div> {snippet.contents || 'no snippet'} </div>
                            </div>
                          )
                        )}
                      </Panel>
                    ))
                  }
                </Col>
              </Row>
            ))
          }
      </Grid>
    )
  }
}

const s2p = state => ({ snippets: state.snippets.weekSnippets })
const d2p = dispatch => bindActionCreators({ loadSnippets }, dispatch)
export default connect(s2p, d2p)(WeekPage)
