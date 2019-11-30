import React, {Component} from 'react'
import {
    Row,
    Col,
    Card
} from 'antd'
import AccountTable from './AccountTable'
// import DebitTable from './DebitTable'
import BreadcrumbCustom from '../BreadcrumbCustom'

import {
    IAccountParam,
    getFinanceAccount,
    setFinanceAccount
} from '../../axios'

type ViewAccountsProps = {
}
interface ViewAccountsState {
    count: number
    data: Array<IAccountParam>
}

class ViewAccounts extends Component<ViewAccountsProps, ViewAccountsState> {
    state = {
        count: 0,
        data: []
    }
    componentDidMount() {
        this.setState({ count: this.state.count+1 })
    }
    componentDidUpdate(prevProps: ViewAccountsProps, prevState: ViewAccountsState) {
        if (this.state.count !== prevState.count) {
            this.getData()
        }
    }
    getData() {
        getFinanceAccount().then(data =>{
            if (!data) return
            this.setState({ data: data })
        })
    }
    onChange(data: IAccountParam) {
        setFinanceAccount('chg', data).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({ count: this.state.count+1 })
            }
        })
    }
    onDelete(data: IAccountParam) {
        setFinanceAccount('del', data).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({ count: this.state.count+1 })
            }
        })
    }
    onCreate(data: IAccountParam) {
        setFinanceAccount('add', data).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({ count: this.state.count+1 })
            }
        })
    }
    render() {
        return (
            <div className="bill">
                <BreadcrumbCustom first="财务" second="账户" />
                <Row>
                    <Col className="bill-row">
                        <Card bordered={false}>
                            <AccountTable
                                data={this.state.data}
                                onChange={ (v) => this.onChange(v) }
                                onDelete={ (v) => this.onDelete(v) }
                                onCreate={ (v) => this.onCreate(v) }
                            />
                        </Card>
                    </Col>
                </Row>
                {/* <Row>
                    <Col span={12}>
                        <Card bordered={false}>
                            <DebitTable
                                title = "借入"
                                data = {this.state.borrow}
                                onChange = {v => this.onChange(null, null, v)}
                                onDelete = {v => this.onDelete(null, null, v)}
                                onCreate = {v => this.onCreate(null, null, v)}
                            />
                        </Card>
                    </Col>
                    <Col span={12}>
                        <Card bordered={false}>
                            <DebitTable
                                title = "借出"
                                data = {this.state.lend}
                                onChange = {v => this.onChange(null, v, null)}
                                onDelete = {v => this.onDelete(null, v, null)}
                                onCreate = {v => this.onCreate(null, v, null)}
                            />
                        </Card>
                    </Col>
                </Row> */}
            </div>
        )
    }
}

export default ViewAccounts