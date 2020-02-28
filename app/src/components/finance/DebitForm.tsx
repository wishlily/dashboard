import React, {Component} from 'react';
import {
    Button,
    Modal,
    Form,
    DatePicker,
    InputNumber,
    Input,
    Select
} from 'antd';
import moment from 'moment';
import { FormComponentProps } from 'antd/lib/form';
import { IAccountParam, minVaildTS } from '../../axios';

const FormItem = Form.Item;
const { TextArea } = Input;
const { Option } = Select;

const timeFmt = 'YYYY-MM-DD HH:mm:ss';

interface IDebitForm {
    id: string;
    type: string;
    member: string;
    amount: number;
    account: string;
    note: string;
    deadline: moment.Moment;
}

type BasicFormsProps = {
    visible: boolean;
    onCancel: () => void;
    onCreate: () => void;
    account: Array<string>;
    data: IAccountParam;
} & FormComponentProps;

interface BasicFormsState {
    data: IAccountParam
}

class BasicForms extends Component<BasicFormsProps, BasicFormsState> {
    state = {
        data: {
            time: moment().format(timeFmt),
            id: '',
            type: '',
            amount: 0,
            member: '',
            account: '',
            note: '',
            deadline: moment(0).format(timeFmt)
        }
    };
    componentDidUpdate(prevProps: BasicFormsProps) {
        if (prevProps.data && prevProps.data !== this.state.data) {
            this.setState({ data: prevProps.data });
            var dealine = "1970-01-01 08:00:00"
            if (prevProps.data.deadline) dealine = prevProps.data.deadline
            prevProps.form.setFieldsValue({
                id: prevProps.data.id,
                type: prevProps.data.type,
                member: prevProps.data.member,
                deadline: moment(dealine, timeFmt),
                accout: prevProps.data.account,
                note: prevProps.data.note,
                amount: prevProps.data.amount,
            });
        }
    }
    render() {
        const { visible, onCancel, onCreate, form, account } = this.props;
        const { getFieldDecorator } = form;

        var accountData: Array<string> = [];
        if (account instanceof Array) {
            accountData = account
        }
        const accountOptions = accountData.map(val => <Option key={val}>{val}</Option>);
        const formItemLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 16 },
        };
        return (
            <Modal
                title="借贷信息"
                visible={visible}
                okText="完成"
                cancelText="取消"
                onCancel={onCancel}
                onOk={onCreate}
            >
                <Form>
                    <FormItem {...formItemLayout} label="姓名" colon={false}>
                        {getFieldDecorator('member', {
                            initialValue: this.state.data.member,
                            rules: [{required: true, message: '请输入姓名'}],
                        })(
                            <Input style={{ width: '100%' }} />
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="金额" colon={false}>
                        {getFieldDecorator('amount', { initialValue: this.state.data.amount })(
                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="账户" colon={false}>
                        {getFieldDecorator('account', { initialValue: this.state.data.account })(
                            <Select showSearch style={{ width: '100%' }} >
                                {accountOptions}
                            </Select>
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="截至日期" colon={false}>
                        {getFieldDecorator('deadline', { initialValue: moment(this.state.data.deadline) })(
                            <DatePicker style={{ width: '100%' }} showTime format={timeFmt} />
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="备注" colon={false}>
                        {getFieldDecorator('note', { initialValue: this.state.data.note })(
                            <TextArea
                                style={{ width: '100%' }}
                                autosize
                            />
                        )}
                    </FormItem>
                </Form>
            </Modal>
        )
    }
}

const CollectionCreateForm: any = Form.create()(BasicForms);

type DebitFormProps = {
    title: string;
    onClick?: () => void;
    onCreate?: (v: IAccountParam) => void;
    account?: Array<string>;
    data?: IAccountParam;
}

class DebitForm extends Component<DebitFormProps> {
    state = {
        visible: false,
    };
    form: any;
    showModal = () => {
        if (this.props.onClick) this.props.onClick();
        this.setState({ visible: true });
    };
    handleCancel = () => {
        this.setState({ visible: false });
    };
    handleCreate = () => {
        const form = this.form;
        form.validateFields((err: any, values: IDebitForm) => {
            if (err) {
                return;
            }
            const data: IAccountParam = {
                time: moment().format(timeFmt),
                id: values.id,
                type: values.type,
                member: values.member,
                amount: values.amount,
                account: values.account,
                note: values.note
            }
            if (values.deadline) {
                if (moment(values.deadline).unix() > minVaildTS()) {
                    data.deadline = values.deadline.format(timeFmt)
                }
            }
            if (this.props.onCreate) this.props.onCreate(data);
            form.resetFields();
            this.setState({ visible: false });
        });
    };
    saveFormRef = (form: any) => {
        this.form = form;
    };
    render() {
        return (
            <span>
                <Button type="primary" onClick={this.showModal}>{this.props.title}</Button>
                <CollectionCreateForm
                    ref={this.saveFormRef}
                    visible={this.state.visible}
                    onCancel={this.handleCancel}
                    onCreate={this.handleCreate}
                    account={this.props.account}
                    data={this.props.data}
                />
            </span>
        );
    }
}

export default DebitForm;