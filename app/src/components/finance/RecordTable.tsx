import React, {Component} from 'react';
import {
    Table,
    Button,
    DatePicker,
    Row,
    Col
} from 'antd';
import moment from 'moment';
import { IRecordParam } from '../../axios';
import RecordForm from './RecordForm';

import {
    getFinanceRecord,
    setFinanceRecord,
    getFinanceAccount
} from '../../axios';

const { MonthPicker } = DatePicker;
const timeFmt = 'YYYY-MM-DD HH:mm:ss';

interface ITable {
    type: string
    time: string
    amount: number
    note: string
}

type RecordTableProps = {
}
interface RecordTableState {
    count: number;
    time: string;
    data: Array<IRecordParam>;
    table: Array<ITable>;
    account: Array<string>;
}

class RecordTable extends Component<RecordTableProps, RecordTableState> {
    constructor(props: any) {
        super(props);
        this.state = {
            count: 0,
            time: moment().format('YYYY-MM'),
            data: [],
            table: [],
            account: [],
        }
        this.columns = [{
            title: '交易类型',
            dataIndex: 'type',
            width: '10%',
        }, {
            title: '日期',
            dataIndex: 'time',
            width: '15%',
        }, {
            title: '金额',
            dataIndex: 'amount',
            width: '10%',
        }, {
            title: '备注',
            dataIndex: 'note',
            width: '45%',
        }, {
            title: '操作',
            dataIndex: 'operation',
            render: (text: any, record: any, index: number) => {
                if (record.type !== '修正') {
                    return (
                        <div className="editable-row-operations">
                            <RecordForm
                                title = "编辑"
                                data = {this.state.data[index]}
                                account = {this.state.account}
                                onCreate = {values => this.onChange(values, index)}
                            />
                            <Button type="primary" onClick={() => this.onDelete(index)}>删除</Button>
                        </div>
                    );
                }
                return (
                    <div className="editable-row-operations">
                        <Button type="primary" onClick={() => this.onDelete(index)}>删除</Button>
                    </div>
                )
            },
        }];
    }
    componentDidMount() {
        getFinanceAccount("list").then(data =>{
            this.setState({account: data})
        })
        this.getData(this.state.time)
    }
    componentDidUpdate(prevProps: RecordTableProps, prevState: RecordTableState) {
        if (this.state.count !== prevState.count || this.state.time !== prevState.time) {
            this.getData(this.state.time)
        }
    }
    columns: any;
    convert(data: IRecordParam) {
        var tabel: ITable = {
            type: data.type,
            time: data.time,
            amount: data.amount,
            note: ""
        }
        if (data.note) tabel.note = data.note
        return tabel
    }
    getData(tmonth: string) {
        const time = moment(tmonth + '-01 00:00:00', timeFmt)
        const t_start = time.format(timeFmt)
        const t_end = time.add(1, 'months').subtract(1, 'seconds').format(timeFmt)
        getFinanceRecord(t_start, t_end).then(result => {
            if (!result) return
            var raw: Array<IRecordParam> = result
            var data: Array<IRecordParam> = []
            var table: Array<ITable> = []
            raw.forEach(record => {
                if (record.type !== '修正') {
                    table.push(this.convert(record))
                    data.push(record)
                }
            });
            this.setState({ data: data, table: table})
        });
    }
    onChange(values: IRecordParam, index: number) {
        const datas = this.state.data
        const data = datas[index]
        values.uuid = data.uuid
        datas[index] = values
        // console.log(data, values)
        setFinanceRecord('chg', values).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({count: this.state.count+1})
            }
        })
    }
    onDelete(index: number) {
        const datas = this.state.data
        const data = datas[index]
        // console.log("delete: ", data);
        setFinanceRecord('del', data).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({count: this.state.count+1})
            }
        })
    }
    onCreate(values: IRecordParam) {
        // console.log("create: ", values)
        setFinanceRecord('add', values).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({count: this.state.count+1})
            }
        })
    }
    onChangeTime(date: any, dateString: string) {
        this.setState({ time: dateString });
    }
    render() {
        return (
            <div>
                <div style={{ marginBottom: 16 }}>
                    <Row gutter={16}>
                        <Col span={2}>
                            <MonthPicker
                                defaultValue={moment(this.state.time, 'YYYY-MM')}
                                onChange={(date, dateString) => this.onChangeTime(date, dateString)}
                                placeholder="请选择月份"
                                style={{ width: 100 }}
                            />
                        </Col>
                        <Col span={2}>
                            <RecordForm
                                title = "新建"
                                account = {this.state.account}
                                onCreate = {values => this.onCreate(values)}
                            />
                        </Col>
                    </Row>
                </div>
                <Table bordered dataSource={this.state.table} columns={this.columns} />
            </div>
        );
    }
}

export default RecordTable;