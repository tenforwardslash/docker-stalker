import {Component} from "react";
import React from "react";
import axios from "axios";
import { Redirect } from 'react-router-dom';

import Constants from "../Constants";
import './ContainerTable.scss';
import '../Utils/Common.scss'

class ContainerTable extends Component {
    constructor(props) {
        super(props);
        /*
        * table will have rows:
        *   IMAGE, STATUS, CREATED, NAME, {STATE, containerID WILL BE NOT VISIBLE}
        * */
        this.state = {
            clickedContainerId: null,
            data : [
                {name:"/clever_lichterman",image:"alpine:latest",created:1544398244,status:"Up 36 minutes",state:"running",ports:null,mounts:null,envVars:["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],networks:["bridge"],containerId:"delta"},
                {name:"/oingo_boingo",image:"alpine:latest",created:1544398244,status:"Up 12 minutes",state:"running",ports:null,mounts:null,envVars:["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],networks:["bridge"],containerId:"gamma"}
            ]
        };
        this.renderItem = this.renderItem.bind(this);
    }
    componentDidMount() {
        let self = this;
        //todo: add to header the password
        axios.get(Constants.API_BASE + "/containers").then(function(response) {
            switch (response.status) {
                case 200:
                    console.debug("successful container retrieval", response.data);
                    self.setState({data: response.data});
                    break;
                case 401:
                    console.error("unauthorized");
                    break;
                default:
                    console.error("cannot handle status code", response);
                    break;

            }
        }).catch(function(error) {
            console.error("error occurred!!!", error.toString())
        })
    }
    handleRowClick(containerId) {
        //redirect to /container/{containerId}/detail
        console.log(containerId);
        this.setState({clickedContainerId: containerId});
    }

    renderItem(item) {
        let createdTime = new Date(0);
        createdTime.setUTCSeconds(item.created);
        const clickCallback = () => this.handleRowClick(item.containerId);
        return (
            [<tr className="row" onClick={clickCallback} key={item.containerId}><td>{item.image}</td><td>{item.status}</td><td>{createdTime.toDateString()}</td><td>{item.name}</td></tr>]
        );
    }

    render() {
        if (this.state.clickedContainerId) {
            let url = `/container/${this.state.clickedContainerId}`;
            return <Redirect push to={url} />
        }
        return <Table rows={this.state.data} renderItem={this.renderItem}/>

        // let allItemRows = [];
        // allItemRows.push((<tr key={"row-data-header"}>
        //     <th>Image</th>
        //     <th>Status</th>
        //     <th>Created</th>
        //     <th>Name</th>
        // </tr>));
        // console.log('rows1', allItemRows);
        // this.state.data.forEach(item => {
        //     const perItemRows = this.renderItem(item);
        //     allItemRows = allItemRows.concat(perItemRows);
        // });
        // console.log('rows2', allItemRows);
        // return (
        //     <table><tbody>{allItemRows}</tbody></table>
        // );
    }
}

const Table = (props) => {
    let allItemRows = [];
    allItemRows.push((<tr className="stalker-bg" key={"row-data-header"}>
        <th>Image</th>
        <th>Status</th>
        <th>Created</th>
        <th>Name</th>
    </tr>));
    props.rows.forEach(item => {
        const perItemRows = props.renderItem(item);
        allItemRows = allItemRows.concat(perItemRows);
    });
    return (
        <div className="container-table"><table><tbody>{allItemRows}</tbody></table></div>
    );
};

export default ContainerTable;
