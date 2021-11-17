import React, { useCallback } from "react";
import { PlaidLink, PlaidLinkOnSuccess } from "react-plaid-link";
import axios from "axios";

interface ILinkProps {
}

interface ILinkState {
  link_token: string;
}

class Link extends React.Component<ILinkProps, ILinkState>{

    constructor(props: ILinkProps){
        super(props);
        this.state = {
            link_token: "" 
        };
    }

    componentDidMount(){
        axios.get("http://localhost:3000/linktoken")
            .then((resp => {
                this.setState({
                    link_token: resp.data
                })
            }), error =>{
                console.log(error);
            }
        );

    }

    render(){
        var { link_token } = this.state;
        return link_token === "" ? (
            // insert your loading animation here
            <div className="loader">Loading</div>
          ) : (
            <PlaidLink
              token={link_token}
              onSuccess={this.onSuccess}
              // onExit={...}
              // onEvent={...}
            >
              Connect a bank account
            </PlaidLink>
        );

    }

    onSuccess = (public_token: any, metadata: any) => {
        console.log(public_token)
        axios({
            method: 'POST',
            url: 'http://localhost:3000/accesstoken',
            data: {
                public_token: public_token,
            },
            headers:{
                "Access-Control-Allow-Origin": "*"
            }
        }).then((resp => {
            console.log(resp)
            axios.get(`http://localhost:3000/transactions/${resp.data}`)

        }), (error => {
          console.log(error);
        })
      );
    }
}

export default Link;

