import React, { Component } from 'react';
import Auth from './Components/Auth.js';
import ContainerTable from './Components/ContainerTable'
import ContainerDetail from './Components/ContainerDetail'
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';
import './App.css';


class App extends Component {
  render() {
    return (
        <Router>
            <div className="App">
                <Route exact path='/' render={()=><Auth />} />
                <Route exact path='/containers' render={()=><ContainerTable/>} />
                <Route exact path="/container/:containerId" component={ContainerDetail} />
            </div>
        </Router>
    );
  }
}

export default App;