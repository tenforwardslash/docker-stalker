import React, { Component } from 'react';
import Constants from "../Constants";

const errorHandler = (WrappedComponent, axios) => {
    return class ErrorHandler extends Component {
        state = {
            error: null
        };

        componentWillMount() {
            // Set axios interceptors
            this.requestInterceptor = axios.interceptors.request.use(req => {
                req.headers.Authorization = localStorage.getItem(Constants.TOKEN_KEY);
                this.setState({ error: null });
                return req;
            });


            this.responseInterceptor = axios.interceptors.response.use(
                res => { return res; },
                error => {
                    this.setState({ error });
                    if (error.response
                        && error.response.status === 401) {
                        this.props.history.push('/');
                        localStorage.removeItem(Constants.TOKEN_KEY);
                    }
                    return error.response;
                }
            );
        }

        componentWillUnmount() {
            // Remove handlers, so Garbage Collector will get rid of if WrappedComponent will be removed
            axios.interceptors.request.eject(this.requestInterceptor);
            axios.interceptors.response.eject(this.responseInterceptor);
        }

        render() {
            let renderSection = this.state.error ? <div style={{"color": "red", "fontSize": "20px"}}>{ this.state.error.response ? String(this.state.error.response.data) : String(this.state.error)}</div> : <WrappedComponent {...this.props} />
            return renderSection;
        }
    };
};

export default errorHandler;