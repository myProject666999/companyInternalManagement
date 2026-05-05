import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, message, Space, Popconfirm, Card, Descriptions, Divider, Tag, Badge, Tabs, InputNumber, DatePicker, Radio } from 'antd';
import { PlusOutlined, DeleteOutlined, EyeOutlined, MailOutlined, InboxOutlined, SendOutlined, SearchOutlined } from '@ant-design/icons';
import api from '../utils/axios';
import { Message, User } from '../types';
import dayjs from 'dayjs';

const { Option } = Select;
const { TextArea } = Input;
const { TabPane } = Tabs;
const { RangePicker } = DatePicker;

const Messages: React.FC = () => {
  const [inboxMessages, setInboxMessages] = useState<Message[]>([]);
  const [outboxMessages, setOutboxMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [replyModalVisible, setReplyModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [viewingMessage, setViewingMessage] = useState<Message | null>(null);
  const [replyingMessage, setReplyingMessage] = useState<Message | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [form] = Form.useForm();
  const [replyForm] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [activeTab, setActiveTab] = useState('inbox');

  const fetchInbox = async (params?: any) => {
    setLoading(true);
    try {
      const response = await api.get('/messages/inbox', { params });
      setInboxMessages(response.data);
    } catch (error) {
      message.error('获取收件箱失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchOutbox = async () => {
    setLoading(true);
    try {
      const response = await api.get('/messages/outbox');
      setOutboxMessages(response.data);
    } catch (error) {
      message.error('获取发件箱失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchUsers = async () => {
    try {
      const response = await api.get('/all-users');
      setUsers(response.data);
    } catch (error) {
      message.error('获取用户列表失败');
    }
  };

  useEffect(() => {
    fetchInbox();
    fetchOutbox();
    fetchUsers();
  }, []);

  const handleTabChange = (key: string) => {
    setActiveTab(key);
  };

  const handleSearch = (values: any) => {
    if (activeTab === 'inbox') {
      fetchInbox(values);
    }
  };

  const handleSend = () => {
    form.resetFields();
    setModalVisible(true);
  };

  const handleView = async (record: Message, isInbox: boolean) => {
    try {
      const response = await api.get(`/messages/${record.id}`);
      const msg = response.data;
      setViewingMessage(msg);
      setDetailVisible(true);
      
      if (isInbox && !msg.is_read) {
        fetchInbox();
      }
    } catch (error) {
      message.error('获取消息详情失败');
    }
  };

  const handleReply = (record: Message) => {
    setReplyingMessage(record);
    replyForm.resetFields();
    replyForm.setFieldsValue({
      subject: `回复: ${record.subject}`,
    });
    setReplyModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/messages/${id}`);
      message.success('删除成功');
      fetchInbox();
      fetchOutbox();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSendSubmit = async (values: any) => {
    try {
      await api.post('/messages', values);
      message.success('发送成功');
      setModalVisible(false);
      fetchOutbox();
    } catch (error: any) {
      message.error(error.response?.data?.error || '发送失败');
    }
  };

  const handleReplySubmit = async (values: any) => {
    if (!replyingMessage) return;
    
    try {
      await api.post(`/messages/${replyingMessage.id}/reply`, {
        content: values.content,
      });
      message.success('回复成功');
      setReplyModalVisible(false);
      fetchInbox();
      fetchOutbox();
    } catch (error: any) {
      message.error(error.response?.data?.error || '回复失败');
    }
  };

  const inboxColumns = [
    {
      title: '状态',
      dataIndex: 'is_read',
      key: 'is_read',
      width: 80,
      render: (isRead: boolean) => (
        <Badge dot={!isRead} offset={[5, 0]}>
          <Tag color={isRead ? 'default' : 'blue'}>
            {isRead ? '已读' : '未读'}
          </Tag>
        </Badge>
      ),
    },
    {
      title: '发件人',
      dataIndex: ['sender', 'name'],
      key: 'sender',
      width: 100,
    },
    {
      title: '主题',
      dataIndex: 'subject',
      key: 'subject',
      ellipsis: true,
    },
    {
      title: '内容摘要',
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
      render: (content: string) => content?.slice(0, 50) + (content?.length > 50 ? '...' : ''),
    },
    {
      title: '发送时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 220,
      render: (_: any, record: Message) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleView(record, true)}
          >
            查看
          </Button>
          <Button
            type="link"
            onClick={() => handleReply(record)}
          >
            回复
          </Button>
          <Popconfirm
            title="确定要删除这条消息吗？"
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

  const outboxColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '收件人',
      dataIndex: ['receiver', 'name'],
      key: 'receiver',
      width: 100,
    },
    {
      title: '主题',
      dataIndex: 'subject',
      key: 'subject',
      ellipsis: true,
    },
    {
      title: '内容摘要',
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
      render: (content: string) => content?.slice(0, 50) + (content?.length > 50 ? '...' : ''),
    },
    {
      title: '发送时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: any, record: Message) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleView(record, false)}
          >
            查看
          </Button>
          <Popconfirm
            title="确定要删除这条消息吗？"
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
          <Form.Item name="sender_name" label="发件人">
            <Input placeholder="请输入发件人姓名" style={{ width: 150 }} />
          </Form.Item>
          <Form.Item name="subject" label="主题">
            <Input placeholder="请输入主题" style={{ width: 150 }} />
          </Form.Item>
          <Form.Item name="is_read" label="状态">
            <Select placeholder="请选择状态" style={{ width: 120 }} allowClear>
              <Option value={false}>未读</Option>
              <Option value={true}>已读</Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
                搜索
              </Button>
              <Button onClick={() => { searchForm.resetFields(); fetchInbox(); fetchOutbox(); }}>
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>

      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<SendOutlined />} onClick={handleSend}>
          发送消息
        </Button>
      </div>

      <Tabs activeKey={activeTab} onChange={handleTabChange}>
        <TabPane
          tab={<span><InboxOutlined /> 收件箱</span>}
          key="inbox"
        >
          <Table
            columns={inboxColumns}
            dataSource={inboxMessages}
            rowKey="id"
            loading={loading}
            bordered
          />
        </TabPane>
        <TabPane
          tab={<span><MailOutlined /> 发件箱</span>}
          key="outbox"
        >
          <Table
            columns={outboxColumns}
            dataSource={outboxMessages}
            rowKey="id"
            loading={loading}
            bordered
          />
        </TabPane>
      </Tabs>

      <Modal
        title="发送消息"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSendSubmit}
        >
          <Form.Item
            name="receiver_id"
            label="收件人"
            rules={[{ required: true, message: '请选择收件人' }]}
          >
            <Select placeholder="请选择收件人" showSearch optionFilterProp="children">
              {users.map((user) => (
                <Option key={user.id} value={user.id}>
                  {user.name} ({user.username})
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="subject"
            label="主题"
            rules={[{ required: true, message: '请输入主题' }]}
          >
            <Input placeholder="请输入主题" />
          </Form.Item>
          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入内容' }]}
          >
            <TextArea placeholder="请输入内容" rows={8} showCount maxLength={5000} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                发送
              </Button>
              <Button onClick={() => setModalVisible(false)}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="回复消息"
        open={replyModalVisible}
        onCancel={() => setReplyModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form
          form={replyForm}
          layout="vertical"
          onFinish={handleReplySubmit}
        >
          <Form.Item
            name="subject"
            label="主题"
          >
            <Input disabled />
          </Form.Item>
          {replyingMessage && (
            <div style={{ marginBottom: 16, padding: 12, backgroundColor: '#f5f5f5', borderRadius: 4 }}>
              <p style={{ marginBottom: 8, fontWeight: 'bold' }}>原消息内容：</p>
              <p style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{replyingMessage.content}</p>
            </div>
          )}
          <Form.Item
            name="content"
            label="回复内容"
            rules={[{ required: true, message: '请输入回复内容' }]}
          >
            <TextArea placeholder="请输入回复内容" rows={6} showCount maxLength={5000} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                回复
              </Button>
              <Button onClick={() => setReplyModalVisible(false)}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="消息详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          viewingMessage && activeTab === 'inbox' && (
            <Button key="reply" type="primary" onClick={() => { setDetailVisible(false); handleReply(viewingMessage); }}>
              回复
            </Button>
          ),
          <Button key="close" onClick={() => setDetailVisible(false)}>
            关闭
          </Button>,
        ]}
        width={700}
      >
        {viewingMessage && (
          <div>
            <Descriptions bordered size="small" column={2}>
              <Descriptions.Item label="主题" span={2}>
                {viewingMessage.subject}
              </Descriptions.Item>
              <Descriptions.Item label="发件人">
                {viewingMessage.sender?.name}
              </Descriptions.Item>
              <Descriptions.Item label="收件人">
                {viewingMessage.receiver?.name}
              </Descriptions.Item>
              <Descriptions.Item label="发送时间" span={2}>
                {viewingMessage.created_at}
              </Descriptions.Item>
            </Descriptions>
            <Divider />
            <div style={{ lineHeight: 1.8, whiteSpace: 'pre-wrap', minHeight: 100 }}>
              {viewingMessage.content}
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default Messages;
