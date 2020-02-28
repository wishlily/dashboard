import React, {Component} from 'react'
import {
    Row,
    Col,
    Card
} from 'antd'
import AccountTable from './AccountTable'
import DebitTable from './DebitTable'
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
    list: Array<string>
    account: Array<IAccountParam>
    borrow: Array<IAccountParam>
    lend: Array<IAccountParam>
}

class ViewAccounts extends Component<ViewAccountsProps, ViewAccountsState> {
    state = {
        count: 0,
        list: [],
        account: [],
        borrow: [],
        lend: []
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
        getFinanceAccount().then(result =>{
            if (!result) return
            var data :Array<IAccountParam> = result
            var account :Array<IAccountParam> = []
            var lend :Array<IAccountParam> = []
            var borrow :Array<IAccountParam> = []
            var list :Array<string> = []
            data.forEach(acct => {
                if (acct.type === 'B') {
                    borrow.push(acct)
                } else if (acct.type === 'L') {
                    lend.push(acct)
                } else {
                    account.push(acct)
                    list.push(acct.id)
                }
            })
            // console.log(data, account, lend, borrow, list)
            this.setState({ account: account, lend: lend, borrow: borrow, list: list })
        })
    }
    onChange(data: IAccountParam, type?: string) {
        switch (type) {
            case 'B':
            case 'L':
                data.type = type
                break;
            default:
                break;
        }
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
    onCreate(data: IAccountParam, type?: string) {
        switch (type) {
            case 'B':
            case 'L':
                data.type = type
                break;
            default:
                break;
        }
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
                                data={this.state.account}
                                onChange={ (v) => this.onChange(v) }
                                onDelete={ (v) => this.onDelete(v) }
                                onCreate={ (v) => this.onCreate(v) }
                            />
                        </Card>
                    </Col>
                </Row>
                <Row>
                    <Col span={12}>
                        <Card bordered={false}>
                            <DebitTable
                                title = "借入"
                                data = {this.state.borrow}
                                account = {this.state.list}
                                onChange={ (v) => this.onChange(v, 'B') }
                                onDelete={ (v) => this.onDelete(v) }
                                onCreate={ (v) => this.onCreate(v, 'B') }
                            />
                        </Card>
                    </Col>
                    <Col span={12}>
                        <Card bordered={false}>
                            <DebitTable
                                title = "借出"
                                data = {this.state.lend}
                                account = {this.state.list}
                                onChange={ (v) => this.onChange(v, 'L') }
                                onDelete={ (v) => this.onDelete(v) }
                                onCreate={ (v) => this.onCreate(v, 'L') }
                            />
                        </Card>
                    </Col>
                </Row>
            </div>
        )
    }
}

export default ViewAccounts