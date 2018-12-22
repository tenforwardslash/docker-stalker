import React, { Component } from 'react';
import Auth from './Components/Auth.js';
import ContainerTable from './Components/ContainerTable'
import ContainerDetail from './Components/ContainerDetail'
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';
import './App.css';
import './Utils/Common.scss'

class App extends Component {
  render() {
    return (
        <div>
            <Router>
                <div className="App">
                    <Link to="/containers" style={{"textDecoration": "none"}}><h1 className="App-header stalker-color">docker stalker</h1></Link><hr className="rule"/>
                    <Route exact path='/' render={()=><Auth />} />
                    <Route exact path='/containers' render={(props)=><ContainerTable {...props}/>} />
                    <Route exact path="/container/:containerId" component={ContainerDetail} />
                </div>
            </Router>
        </div>
    );
  }
}

export default App;
