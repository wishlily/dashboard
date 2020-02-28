import React, { Component } from 'react';
import {
    Table,
    Button,
    Row,
    Col,
    Icon
} from 'antd';
// import moment from 'moment';

import DebitForm from './DebitForm';
import { IAccountParam } from '../../axios';

// const timeFmt = 'YYYY-MM-DD HH:mm:ss';

interface IDebitTable {
    key?: string;
    member?: string;
    amount: number;
    data: Array<IAccountParam>;
}

type DebitTableProps = {
    onChange?: (v: IAccountParam) => void;
    onDelete?: (v: IAccountParam) => void;
    onCreate?: (v: IAccountParam) => void;
    title: string;
    account: Array<string>;
    data: Array<IAccountParam>;
}

interface DebitTableState {
    amount: number;
    raw: Array<IAccountParam>;
    data: Array<IDebitTable>;
}

class DebitTable extends Component<DebitTableProps, DebitTableState> {
    constructor(props: any) {
        super(props);
        this.state = {
            raw: [],
            amount: 0,
            data: [],
        };
        this.columns = [{
            title: '姓名',
            dataIndex: 'member',
        }, {
            title: '金额',
            dataIndex: 'amount',
        }];
    }
    static getDerivedStateFromProps(prevProps: DebitTableProps, prevState: DebitTableState) {
        const {data} = prevProps
        if (data && data !== prevState.raw) {
            var amount = 0
            const groups: Array<string|undefined> = []
            data.forEach(item => {
                amount = +Number(amount + item.amount).toFixed(2)
                if (groups.indexOf(item.member) < 0) {
                    groups.push(item.member)
                }
            })
            const items: Array<IDebitTable> = []
            groups.forEach(member => {
                const item: IDebitTable = {key: member, member: member, amount: 0, data: []}
                data.forEach(v => {
                    if (v.member === member) {
                        item.amount = +Number(item.amount + v.amount).toFixed(2)
                        var t: any = v
                        t.key = v.id
                        item.data.push(t)
                    }
                })
                items.push(item)
            })
            return {
                raw: data,
                data: items,
                amount: amount
            };
        }
    }
    columns: any;
    expandedRowRender(record: IDebitTable, account: Array<string>) {
        const columns = [
            {
                title: '金额',
                dataIndex: 'amount'
            }, {
                title: '备注',
                dataIndex: 'note'
            }, {
                title: '还款日期',
                dataIndex: 'deadline'
            }, {
                title: '操作',
                dataIndex: 'operation',
                render: (text: any, record: any, index: number) => {
                    return (
                        <div className="editable-row-operations">
                            <DebitForm
                                title = "+"
                                data = {record}
                                account = {account}
                                onCreate = {values => this.onChange(values)}
                            />
                            <Button type="primary" onClick={() => this.onDelete(record)}>-</Button>
                        </div>
                    );
                },
            }];
        return (
            <Table
                columns={columns}
                dataSource={record.data}
                pagination={false}
            />
        );
    }
    onChange(data: IAccountParam) {
        // same as source
        if (this.props.onChange) this.props.onChange(data);
        // console.log("change: ", values)
    }
    onDelete(data: IAccountParam) {
        // same as source
        if (this.props.onDelete) this.props.onDelete(data);
        // console.log("delete: ", values)
    }
    onCreate(data: IAccountParam) {
        // same as source
        if (this.props.onCreate) this.props.onCreate(data);
        // console.log("create: ", values)
    }
    render() {
        return (
            <div>
                <div style={{ marginBottom: 16 }}>
                    <Row gutter={16} align="bottom">
                        <Col span={4}>
                            <DebitForm
                                title = {this.props.title}
                                account = {this.props.account}
                                onCreate = {values => this.onCreate(values)}
                            />
                        </Col>
                        <Col span={4}>
                            <Icon type="heart" style={{ color: '#FF0040' }} /> {this.state.amount}
                        </Col>
                    </Row>
                </div>
                <Table
                    columns={this.columns}
                    dataSource={this.state.data}
                    expandedRowRender={(record) => this.expandedRowRender(record, this.props.account)}
                    pagination={false}
                />
            </div>
        );
    }
}

export default DebitTable;