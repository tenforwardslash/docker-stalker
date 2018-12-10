import {Component} from "react";
import React from "react";
import axios from "axios";

import Constants from "../Constants";

import "./ContainerDetail.css"

class ContainerDetail extends Component {
    constructor(props) {
        super(props);
        this.state = {
            detail: null
        };
        this.restartContainer = this.restartContainer.bind(this);
    }
    componentDidMount() {
        let self = this;
        axios.get(Constants.API_BASE + `/container/${self.props.match.params.containerId}/detail`).then(function (response){
            switch (response.status) {
                case 200:
                    console.log("nailed it", response.data);
                    self.setState({detail: response.data});
                    break;
                case 401:
                    console.log("unauthorized");
                    break;
                default:
                    console.log("unhandleable status", response);
                    break;
            }
        }).catch(function(error){
            console.error("unable to get container detail", error);
        })
    };

    restartContainer() {
        let id = this.state.detail.containerId;
        console.log(id);
        axios.post(Constants.API_BASE + `/container/${id}/restart`).then(function(response) {
            console.log(response);
            switch (response.status) {
                case 200:
                    console.log("successfully restarted");
                    break;
                case 401:
                    console.log("unauthorized");
                    break;
                default:
                    console.log("no idea");
                    break;
            }
        }).catch(function(error) {console.error("unable to restart", error)})
    };

    render() {
        console.log('detail', this.state.detail);
        if (this.state.detail) {
            return <div> <Detail container={this.state.detail} restartContainer={this.restartContainer}/></div>
        }
        return <div>{this.props.match.params.containerId}</div>
    };
}

const Detail = (props) => {
    console.log("CONTAINER!!!!", props.container.envVars);

    let mounts = props.container.mounts.map((mount) =>
        <li>{mount}</li>
    );
    return (
       <div>
           <div>
               <h1>Container {props.container.image}</h1>
               <button onClick={props.restartContainer}>
                   Restart Container
               </button>
               <div className="Section">
                   <h2>Summary</h2>
                   <ul>
                       <li><b>Container Name: </b>{props.container.name}</li>
                       <li><b>Status: </b>{props.container.status}</li>
                       <li><b>Id: </b>{props.container.containerId}</li>
                       {mounts}
                   </ul>
               </div>
               <br/>
               <Networks networks={props.container.networks}/>
               <Envs envVars={props.container.envVars}/>
           </div>
       </div>
   )
};

const Networks = (props) => {
    if (props.networks) {
        let networks = props.networks.map((network) =>
            <li>{network}</li>
        );
        return (
            <div className="Section">
                <h2>Networks</h2>
                <ul>{networks}</ul>
            </div>
        )
    }
};

const Envs = (props) => {
    let envVars = props.envVars.map((env) =>
        <li>{env}</li>
    );
    if (props.envVars) {
        return (
            <div className="Section">
                <h2>Environment Variables</h2>
                <ul>{envVars}</ul>
            </div>
        )
    }
};

export default ContainerDetail;