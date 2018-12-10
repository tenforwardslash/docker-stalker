import React, { Component } from 'react';
import axios from 'axios';
import { Redirect } from 'react-router-dom';

import Constants from '../Constants';
import './Auth.css';

class Auth extends Component {
    constructor(props) {
        super(props);
        this.state = {value: '', secureErr: null, passThroughError: false, redirect: false};

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    componentDidMount() {
        let self = this;
        axios.get(Constants.API_BASE + '/isSecure').then(function (response) {
            self.setState({
                redirect: !response.data.isSecure,
            });
          }).catch(function (error) {
              console.error("unable to get secure state!!", error);
              self.setState({secureErr: error})
          });
    }

    handleChange(event) {
        this.setState({value: event.target.value});
    }

    handleSubmit(event) {
        let self = this;
        axios.post(Constants.API_BASE + "/login", {password: this.state.value}).then(function (response) {
            let token = response.data.token;
            if (response.status === 200) {
                self.setState({redirect: true});
                localStorage.setItem("stalkerToken", token);
            } else if (response.status === 401) {
                console.log("bad login password")
            } else {
                console.log("unexpected status")
            }
        });
        console.log(this.state.value);
        event.preventDefault();
    }
    render() {
        if (this.state.redirect) {
            return <Redirect to='/containers'/>
        } else if (this.state.secureErr && this.state.passThroughError === true) {
            return <Wrapper component={<div>{this.state.secureErr.toString()}</div>}/>
        }
        return (
            <Wrapper component={<PasswordForm handleSubmit={this.handleSubmit} handleChange={this.handleChange} value={this.state.value}/>}/>
        );
    }
}

const Wrapper = (props) => {
    return (
        <div className="outer-div">
            <div className="inner-div">
                {props.component}
            </div>
        </div>
    )
};

const PasswordForm = (props) => {
    return (
        <form onSubmit={props.handleSubmit}>
            <label>
                Password:
                <input type="password" value={props.value} onChange={props.handleChange} />
            </label>
            <input type="submit" value="Submit" />
        </form>
    )
};

export default Auth;