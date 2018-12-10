import {Component} from "react";
import React from "react";
import axios from "axios";

import Constants from "../Constants";

class ContainerDetail extends Component {
    constructor(props) {
        super(props);
        this.state = {
            detail: null
        }
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

    render() {
        console.log('detail', this.state.detail);
        if (this.state.detail) {
            return <div> <Detail container={this.state.detail}/></div>
        }
        return <div>{this.props.match.params.containerId}</div>
    };
}

const Detail = (props) => {
    console.log("CONTAINER!!!!", props.container.envVars);
    const listItems = numbers.map((number) =>
        <li>{number}</li>
    );
    return (
       <div>
           <div>
               <ul>
                   <li>{props.container.containerId}</li>
                   {props.container.envVars.forEach(item => {return <li>{item}</li>})}
               </ul>
           </div>
       </div>
   )
};

export default ContainerDetail;