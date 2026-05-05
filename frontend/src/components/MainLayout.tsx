import React, { useState } from 'react';
import { Layout, Menu, Dropdown, Avatar, Button, Modal, Form, Input, message } from 'antd';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  TeamOutlined,
  UserOutlined,
  NotificationOutlined,
  FileOutlined,
  MessageOutlined,
  CalendarOutlined,
  ScheduleOutlined,
  LogoutOutlined,
  KeyOutlined,
} from '@ant-design/icons';
import { Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Role } from '../types';
import api from '../utils/axios';
import Departments from '../pages/Departments';
import Users from '../pages/Users';
import Information from '../pages/Information';
import Documents from '../pages/Documents';
import Messages from '../pages/Messages';

const { Sider, Header, Content } = Layout;

interface MenuItem {
  key: string;
  icon: React.ReactNode;
  label: string;
  path: string;
  roles: Role[];
}

const menuItems: MenuItem[] = [
  {
    key: 'dashboard',
    icon: <DashboardOutlined />,
    label: '首页',
    path: '/',
    roles: ['general_manager', 'department_manager', 'employee'],
  },
  {
    key: 'departments',
    icon: <TeamOutlined />,
    label: '部门管理',
    path: '/departments',
    roles: ['general_manager'],
  },
  {
    key: 'users',
    icon: <UserOutlined />,
    label: '人事管理',
    path: '/users',
    roles: ['general_manager'],
  },
  {
    key: 'department-users',
    icon: <UserOutlined />,
    label: '部门员工管理',
    path: '/department-users',
    roles: ['department_manager'],
  },
  {
    key: 'tasks',
    icon: <ScheduleOutlined />,
    label: '工作任务管理',
    path: '/tasks',
    roles: ['department_manager'],
  },
  {
    key: 'department-attendance',
    icon: <CalendarOutlined />,
    label: '部门考勤管理',
    path: '/department-attendance',
    roles: ['department_manager'],
  },
  {
    key: 'my-attendance',
    icon: <CalendarOutlined />,
    label: '个人考勤管理',
    path: '/my-attendance',
    roles: ['employee'],
  },
  {
    key: 'work-logs',
    icon: <FileOutlined />,
    label: '个人办公管理',
    path: '/work-logs',
    roles: ['employee'],
  },
  {
    key: 'information',
    icon: <NotificationOutlined />,
    label: '信息发布管理',
    path: '/information',
    roles: ['general_manager', 'department_manager'],
  },
  {
    key: 'company-information',
    icon: <NotificationOutlined />,
    label: '公司信息查询',
    path: '/company-information',
    roles: ['employee'],
  },
  {
    key: 'documents',
    icon: <FileOutlined />,
    label: '共享文档',
    path: '/documents',
    roles: ['general_manager', 'department_manager', 'employee'],
  },
  {
    key: 'messages',
    icon: <MessageOutlined />,
    label: '个人消息管理',
    path: '/messages',
    roles: ['general_manager', 'department_manager', 'employee'],
  },
];

const MainLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [changePasswordModalVisible, setChangePasswordModalVisible] = useState(false);
  const [form] = Form.useForm();
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const filteredMenuItems = menuItems.filter((item) => user && item.roles.includes(user.role));

  const getSelectedKey = () => {
    const item = filteredMenuItems.find((i) => location.pathname === i.path);
    return item ? [item.key] : ['dashboard'];
  };

  const handleMenuClick = ({ key }: { key: string }) => {
    const item = filteredMenuItems.find((i) => i.key === key);
    if (item) {
      navigate(item.path);
    }
  };

  const handleLogout = () => {
    logout();
  };

  const handleChangePassword = async (values: { old_password: string; new_password: string }) => {
    try {
      await api.post('/auth/change-password', values);
      message.success('密码修改成功，请重新登录');
      setChangePasswordModalVisible(false);
      form.resetFields();
      logout();
    } catch (error: any) {
      message.error(error.response?.data?.error || '密码修改失败');
    }
  };

  const userMenuItems = [
    {
      key: 'change-password',
      icon: <KeyOutlined />,
      label: '修改密码',
      onClick: () => setChangePasswordModalVisible(true),
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  const roleNames: Record<Role, string> = {
    general_manager: '总经理',
    department_manager: '部门经理',
    employee: '普通员工',
  };

  return (
    <Layout className="layout-container">
      <Sider trigger={null} collapsible collapsed={collapsed} className="sider">
        <div className="logo">管理系统</div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={getSelectedKey()}
          items={filteredMenuItems.map((item) => ({
            key: item.key,
            icon: item.icon,
            label: item.label,
          }))}
          onClick={handleMenuClick}
        />
      </Sider>
      <Layout className="site-layout">
        <Header className="site-layout-header">
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '16px', width: 64, height: 64 }}
          />
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <span style={{ color: '#666' }}>
              {user?.name} ({roleNames[user?.role || 'employee']})
            </span>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Avatar icon={<UserOutlined />} style={{ cursor: 'pointer' }} />
            </Dropdown>
          </div>
        </Header>
        <Content className="site-layout-content">
          <Routes>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/departments" element={<Departments />} />
            <Route path="/users" element={<Users />} />
            <Route path="/department-users" element={<Users />} />
            <Route path="/tasks" element={<div>工作任务管理页面</div>} />
            <Route path="/department-attendance" element={<div>部门考勤管理页面</div>} />
            <Route path="/my-attendance" element={<div>个人考勤管理页面</div>} />
            <Route path="/work-logs" element={<div>个人办公管理页面</div>} />
            <Route path="/information" element={<Information />} />
            <Route path="/company-information" element={<Information />} />
            <Route path="/documents" element={<Documents />} />
            <Route path="/messages" element={<Messages />} />
          </Routes>
        </Content>
      </Layout>

      <Modal
        title="修改密码"
        open={changePasswordModalVisible}
        onCancel={() => {
          setChangePasswordModalVisible(false);
          form.resetFields();
        }}
        footer={null}
      >
        <Form
          form={form}
          onFinish={handleChangePassword}
          layout="vertical"
        >
          <Form.Item
            name="old_password"
            label="旧密码"
            rules={[{ required: true, message: '请输入旧密码' }]}
          >
            <Input.Password placeholder="请输入旧密码" />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码长度至少6位' },
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认新密码"
            dependencies={['new_password']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password placeholder="请确认新密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              确认修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  );
};

const DashboardPage: React.FC = () => {
  const { user } = useAuth();

  const roleNames: Record<Role, string> = {
    general_manager: '总经理',
    department_manager: '部门经理',
    employee: '普通员工',
  };

  return (
    <div>
      <h2 className="page-title">欢迎使用公司内部管理系统</h2>
      <div style={{ marginBottom: 24, padding: 24, background: '#f5f5f5', borderRadius: 8 }}>
        <h3>当前用户信息</h3>
        <p>姓名: {user?.name}</p>
        <p>用户名: {user?.username}</p>
        <p>角色: {roleNames[user?.role || 'employee']}</p>
        <p>邮箱: {user?.email || '未设置'}</p>
        <p>电话: {user?.phone || '未设置'}</p>
      </div>
      <div>
        <h3>使用说明</h3>
        <ul style={{ lineHeight: 2 }}>
          <li>左侧菜单可以点击展开/收起</li>
          <li>点击右上角头像可以修改密码或退出登录</li>
          <li>根据您的角色不同，菜单权限也会有所不同</li>
        </ul>
      </div>
    </div>
  );
};

export default MainLayout;
