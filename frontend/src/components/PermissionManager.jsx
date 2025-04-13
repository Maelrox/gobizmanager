import React, { useState, useEffect } from 'react';
import { Modal, Select, message } from 'antd';
import axios from 'axios';

const PermissionManager = ({ role, onClose }) => {
  const [users, setUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const response = await axios.get(`/api/users/search?companyId=${role.company_id}`);
      setUsers(response.data);
    } catch (error) {
      message.error('Failed to fetch users');
    }
  };

  const handleAssignUser = async () => {
    if (!selectedUser) {
      message.error('Please select a user');
      return;
    }

    setLoading(true);
    try {
      await axios.post('/api/roles/assign', {
        userId: selectedUser,
        roleId: role.id
      });
      message.success('User assigned successfully');
      onClose();
    } catch (error) {
      message.error('Failed to assign user');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title="Assign User to Role"
      open={true}
      onOk={handleAssignUser}
      onCancel={onClose}
      confirmLoading={loading}
    >
      <Select
        style={{ width: '100%' }}
        placeholder="Select a user"
        onChange={setSelectedUser}
        options={users.map(user => ({
          value: user.id,
          label: user.name
        }))}
      />
    </Modal>
  );
};

export default PermissionManager; 