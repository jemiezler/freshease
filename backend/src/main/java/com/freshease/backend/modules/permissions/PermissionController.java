/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */

package com.freshease.backend.modules.permissions;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import lombok.RequiredArgsConstructor;

/**
 *
 * @author jemiezler
 */
@RestController
@RequestMapping("/permissions")
@RequiredArgsConstructor
public class PermissionController {
    private final PermissionService permissionService;

    @PostMapping()
    public PermissionEntity createPerm(@RequestParam String name,
            @RequestParam(required = false) String resourceType,
            @RequestParam(required = false) String resourceId) {
        return permissionService.createPermission(
                new PermissionEntity(
                        java.util.UUID.randomUUID(),
                        name,
                        resourceType,
                        resourceId,
                        java.time.OffsetDateTime.now(),
                        java.time.OffsetDateTime.now()));
    }
}
