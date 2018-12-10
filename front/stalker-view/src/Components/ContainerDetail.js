import {Component} from "react";
import React from "react";
import axios from "axios";

import Constants from "../Constants";

import "./ContainerDetail.scss"
import "../Utils/Common.scss"

const RestartEnum = Object.freeze({"active":1, "success":2, "failed": 3});


class ContainerDetail extends Component {
    constructor(props) {
        super(props);
        this.state = {
            detail: null,
            restartState: null,

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
                    console.log("unhandle-able status", response);
                    break;
            }
        }).catch(function(error){
            console.error("unable to get container detail", error);
        })
    };

    restartContainer() {
        let self = this;
        let id = this.state.detail.containerId;
        console.log(id);
        this.setState({restartState: RestartEnum.active});
        axios.post(Constants.API_BASE + `/container/${id}/restart`).then(function(response) {
            console.log(response);
            switch (response.status) {
                case 200:
                    console.log("successfully restarted");
                    self.setState({restartState: RestartEnum.success});
                    break;
                case 401:
                    self.setState({restartState: RestartEnum.failed});
                    console.log("unauthorized");
                    break;
                default:
                    self.setState({restartState: RestartEnum.failed});
                    console.log("no idea");
                    break;
            }
        }).catch(function(error) {console.error("unable to restart", error)})
    };

    render() {
        console.log('detail', this.state.detail);
        if (this.state.detail) {
            return <div> <Detail container={this.state.detail} restartState={this.state.restartState} restartContainer={this.restartContainer}/></div>
        }
        return <div>{this.props.match.params.containerId}</div>
    };
}

const Detail = (props) => {
    console.log("CONTAINER!!!!", props.container.envVars);

    let mounts = props.container.mounts.map((mount) =>
        <li>{mount}</li>
    );
    let splitImage = props.container.image.split(":");
    return (
       <div>
           <div>
               <h1 className="Header">container {splitImage[0]}:<b className="stalker-color">{splitImage[1]}</b></h1>
               <Restart restartContainer={props.restartContainer} restartState={props.restartState}/>
               <div className="Section">
                   <h2 className="SectionHeader">Summary</h2>
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

const Restart = (props) => {
    let restart = null;
    switch (props.restartState) {
        case RestartEnum.failed:
            restart = <div><b>Failed to restart</b></div>;
            break;
        case RestartEnum.active:
            restart = <div>Restarting...</div>;
            break;
        case RestartEnum.success:
            restart = <div>Successfully restarted!</div>;
            break;
        default:
            restart = <button className="stalker-button" onClick={props.restartContainer}>Restart Container</button>;
            break;
    }
    return restart
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