import axios from "axios";

const Constants = {
    API_BASE: process.env.REACT_APP_API_SERVER,
    TOKEN_KEY: "stalkerToken",
    API: axios.create({
        baseURL: process.env.REACT_APP_API_SERVER,
        timeout: 10000,
        headers: {
            'Accept': 'application/json'
        }
    })
};

export default Constants;