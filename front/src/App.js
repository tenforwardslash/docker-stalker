import React, { Component } from 'react';
import Auth from './Components/Auth.js';
import ContainerTable from './Components/ContainerTable'
import ContainerDetail from './Components/ContainerDetail'
import { BrowserRouter as Router, Route } from 'react-router-dom';
import './App.css';
import './Utils/Common.scss'

class App extends Component {
  render() {
    return (
        <div>
            <a style={{"textDecoration": "none"}}href="/containers"><h1 className="App-header stalker-color">docker stalker</h1></a><hr className="rule"/>
            <Router>
                <div className="App">
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
