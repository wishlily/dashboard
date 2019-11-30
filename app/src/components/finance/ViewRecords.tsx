import React, {Component} from 'react';
import {
    Row,
    Col,
    Card
} from 'antd';
import RecordTable from './RecordTable';
import BreadcrumbCustom from '../BreadcrumbCustom';

class ViewRecords extends Component {
    render() {
        return (
            <div className="bill">
                <BreadcrumbCustom first="财务" second="流水" />
                <Row gutter={16}>
                    <Col className="bill-row" md={24}>
                        <div className="bill-box">
                            <Card title="" bordered={false}>
                                <RecordTable />
                            </Card>
                        </div>
                    </Col>
                </Row>
            </div>
        )
    }
}

export default ViewRecords;