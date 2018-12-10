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
                   </ul>
               </div>
               <br/>
               <ListSection properties={props.container} section={"networks"} sectionTitle={"Networks"}/>
               <ListSection properties={props.container} section={"envVars"} sectionTitle={"Environment Variables"}/>
               {/* FIXME mounts isn't a simple array, will need its own code
                <Section properties={props.container} section={"mounts"} sectionTitle={"Volume Mounts"}/>
               */}
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

// ListSection is for a section that has a list of strings as its values. each value of the list will be written as an html list item
const ListSection = (props) => {
  if (props.properties[props.section]) {
      let sectionProps = props.properties[props.section];
      if (sectionProps.length === 0) {
          return <div/>
      }
      console.log("SECTION", props.section, sectionProps);
      let sectionList = sectionProps.map((elem) =>
          <li key={"section-item-"+elem}>{elem}</li>
      );
      return (
          <div className="Section">
              <h2 className="SectionHeader">{props.sectionTitle}</h2>
              <ul>{sectionList}</ul>
          </div>
      )
  }
};


export default ContainerDetail;