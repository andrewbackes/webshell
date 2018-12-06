import React, { Component } from 'react';
import Presentation from './presentation';

export default class Terminal extends Component {

    constructor(props) {
        super(props);
        this.state = {
            content: [],
            input: ''
        };
        this.bindInput();
        this.receiveLine = this.receiveLine.bind(this);
        this.initWebsocket();
    }

    receiveLine(line) {
        let lines = this.state.content;
        lines.push(line);
        this.setState({ content: lines });
    }

    bindInput() {
        const this2 = this;
        document.addEventListener('keyup', function (event) {
            if (event.defaultPrevented) {
                return;
            }
            var key = event.key || event.keyCode;
            if (key === 'Escape' || key === 'Esc' || key === 27) {
                this2.setState({ input: '' });
            } else if (key === "Enter" || key === 13) {
                this2.ws.send(this2.state.input);
                this2.setState({ input: '' });
            } else if (key === 'Backspace' || key === 8) {
                this2.setState({ input: this2.state.input.slice(0, -1) });
            } else {
                // invalid: 0-46, 91-105, 112-145
                const k = event.keyCode;
                const invalid = (k < 32) || (k > 32 && k < 47) || (k > 90 && k < 106) || (k > 111 && k < 146);
                if (!invalid) {
                    this2.setState({ input: this2.state.input + key });
                }
            }
        });
    }

    initWebsocket() {
        this.ws = new WebSocket(((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/websocket");
        const this2 = this;
        this.ws.onmessage = function (event) {
            this2.receiveLine(event.data);
        };
        this.ws.onclose = function () {
            this2.receiveLine("Websocket closed.");
        };
        this.ws.onerror = function (event) {
            this2.receiveLine("WebSocket error: " + event);
        };
    }

    render() {
        return (
            <Presentation content={this.state.content} input={this.state.input} />
        );
    }
}

