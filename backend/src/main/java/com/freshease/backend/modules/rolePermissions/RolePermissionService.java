/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */

package com.freshease.backend.modules.rolePermissions;

import java.util.UUID;

import org.springframework.stereotype.Service;

import com.freshease.backend.modules.permissions.PermissionRepository;
import com.freshease.backend.modules.roles.RoleRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;

/**
 *
 * @author jemiezler
 */
@Service
@RequiredArgsConstructor
public class RolePermissionService {
    private final RolePermissionRepository rolePermissionRepository;
    private final RoleRepository roleRepository;
    private final PermissionRepository permissionRepository;

    @Transactional
    public void grantPermissionToRole(UUID roleId, UUID permId) {
        if (rolePermissionRepository.existsByRole_IdAndPermission_Id(roleId, permId))
            return;
        var role = roleRepository.findById(roleId).orElseThrow();
        var perm = permissionRepository.findById(permId).orElseThrow();
        rolePermissionRepository.save(RolePermissionEntity.builder().role(role).permission(perm).build());
    }
}
