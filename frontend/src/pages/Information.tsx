import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, message, Space, Popconfirm, Card, Descriptions, Tag, Divider } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined, SearchOutlined } from '@ant-design/icons';
import api from '../utils/axios';
import { Information } from '../types';

const { Option } = Select;
const { TextArea } = Input;

const InformationPage: React.FC = () => {
  const [informations, setInformations] = useState<Information[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [editingInfo, setEditingInfo] = useState<Information | null>(null);
  const [viewingInfo, setViewingInfo] = useState<Information | null>(null);
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();

  const typeOptions = [
    { value: 'notice', label: '通知公告' },
    { value: 'news', label: '新闻动态' },
    { value: 'policy', label: '规章制度' },
    { value: 'other', label: '其他' },
  ];

  const fetchInformations = async (params?: any) => {
    setLoading(true);
    try {
      const response = await api.get('/information', { params });
      setInformations(response.data);
    } catch (error) {
      message.error('获取信息列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInformations();
  }, []);

  const handleSearch = (values: any) => {
    fetchInformations(values);
  };

  const handleAdd = () => {
    setEditingInfo(null);
    form.resetFields();
    form.setFieldsValue({
      type: 'notice',
      is_public: true,
    });
    setModalVisible(true);
  };

  const handleEdit = (record: Information) => {
    setEditingInfo(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleView = async (record: Information) => {
    try {
      const response = await api.get(`/information/${record.id}`);
      setViewingInfo(response.data);
      setDetailVisible(true);
    } catch (error) {
      message.error('获取信息详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/information/${id}`);
      message.success('删除成功');
      fetchInformations();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSubmit = async (values: any) => {
    try {
      if (editingInfo) {
        await api.put(`/information/${editingInfo.id}`, values);
        message.success('更新成功');
      } else {
        await api.post('/information', values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchInformations();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const getTypeLabel = (type: string) => {
    const option = typeOptions.find((o) => o.value === type);
    return option ? option.label : type;
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 120,
      render: (type: string) => getTypeLabel(type),
    },
    {
      title: '作者',
      dataIndex: ['author', 'name'],
      key: 'author',
      width: 100,
    },
    {
      title: '是否公开',
      dataIndex: 'is_public',
      key: 'is_public',
      width: 100,
      render: (isPublic: boolean) => (
        <Tag color={isPublic ? 'green' : 'orange'}>
          {isPublic ? '公开' : '内部'}
        </Tag>
      ),
    },
    {
      title: '查看次数',
      dataIndex: 'view_count',
      key: 'view_count',
      width: 100,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 220,
      render: (_: any, record: Information) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleView(record)}
          >
            查看
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这条信息吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Form
          form={searchForm}
          layout="inline"
          onFinish={handleSearch}
        >
          <Form.Item name="type" label="类型">
            <Select placeholder="请选择类型" style={{ width: 150 }} allowClear>
              {typeOptions.map((option) => (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="title" label="标题">
            <Input placeholder="请输入标题" style={{ width: 200 }} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
                搜索
              </Button>
              <Button onClick={() => { searchForm.resetFields(); fetchInformations(); }}>
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>

      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          发布信息
        </Button>
      </div>
      
      <Table
        columns={columns}
        dataSource={informations}
        rowKey="id"
        loading={loading}
        bordered
      />

      <Modal
        title={editingInfo ? '编辑信息' : '发布信息'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={700}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input placeholder="请输入标题" />
          </Form.Item>
          <Form.Item
            name="type"
            label="类型"
            rules={[{ required: true, message: '请选择类型' }]}
          >
            <Select placeholder="请选择类型">
              {typeOptions.map((option) => (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入内容' }]}
          >
            <TextArea placeholder="请输入内容" rows={10} showCount maxLength={10000} />
          </Form.Item>
          <Form.Item
            name="is_public"
            label="是否公开"
            valuePropName="checked"
          >
            <Select>
              <Option value={true}>公开</Option>
              <Option value={false}>内部</Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                确定
              </Button>
              <Button onClick={() => setModalVisible(false)}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="信息详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailVisible(false)}>
            关闭
          </Button>,
        ]}
        width={800}
      >
        {viewingInfo && (
          <div>
            <h2 style={{ marginBottom: 16 }}>{viewingInfo.title}</h2>
            <Descriptions bordered size="small" column={3}>
              <Descriptions.Item label="类型">
                {getTypeLabel(viewingInfo.type)}
              </Descriptions.Item>
              <Descriptions.Item label="作者">
                {viewingInfo.author?.name}
              </Descriptions.Item>
              <Descriptions.Item label="查看次数">
                {viewingInfo.view_count}
              </Descriptions.Item>
              <Descriptions.Item label="是否公开">
                {viewingInfo.is_public ? '公开' : '内部'}
              </Descriptions.Item>
              <Descriptions.Item label="发布时间" span={2}>
                {viewingInfo.created_at}
              </Descriptions.Item>
            </Descriptions>
            <Divider />
            <div style={{ lineHeight: 1.8, whiteSpace: 'pre-wrap' }}>
              {viewingInfo.content}
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default InformationPage;
