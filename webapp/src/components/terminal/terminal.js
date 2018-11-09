import React, { Component } from 'react';
import Presentation from './presentation';

export default class Terminal extends Component {
    constructor(props) {
        super(props);
        this.state = {
            content: ['line 1', 'line 2'],
            input: ''
        };
        this.bindInput();
    }

    bindInput() {
        const this2 = this;
        document.addEventListener('keypress', function (event) {
            if (event.defaultPrevented) {
                return;
            }
            var key = event.key || event.keyCode;
            if (key === 'Escape' || key === 'Esc' || key === 27) {
                this2.setState({ input: '' });
            } else if (key === "Enter") {
                // todo: submit
                this2.setState({ input: '' });
            }
            else {
                this2.setState({ input: this2.state.input + key });
            }
        });
    }

    render() {
        return (
            <Presentation content={this.state.content} input={this.state.input} />
        );
    }
}

