package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"watchducker/internal/docker"
	"watchducker/internal/types"
	"watchducker/pkg/logger"
	"watchducker/pkg/utils"
)

// Checker 核心检查器
type Checker struct {
	clientManager *docker.ClientManager
	containerSvc  *docker.ContainerService
	imageSvc      *docker.ImageService
}

// NewChecker 创建新的检查器实例
func NewChecker() (*Checker, error) {
	clientManager, err := docker.NewClientManager()
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端管理器失败: %w", err)
	}

	containerSvc := docker.NewContainerService(clientManager)
	imageSvc := docker.NewImageService(clientManager)

	return &Checker{
		clientManager: clientManager,
		containerSvc:  containerSvc,
		imageSvc:      imageSvc,
	}, nil
}

// CheckByName 根据容器名称检查镜像更新
func (c *Checker) CheckByName(ctx context.Context, containerNames []string) (*types.BatchCheckResult, error) {
	logger.Info("开始根据容器名称检查镜像更新: %v", containerNames)

	// 获取所有指定名称的容器
	containers, err := c.containerSvc.GetByName(ctx, containerNames)
	if err != nil {
		return nil, fmt.Errorf("获取容器失败: %w", err)
	}

	// 使用通用检查逻辑
	return c.checkImages(ctx, containers, utils.CreateCheckCallback())
}

// CheckByLabel 根据标签检查镜像更新
func (c *Checker) CheckByLabel(ctx context.Context, labelKey, labelValue string) (*types.BatchCheckResult, error) {
	logger.Info("开始根据标签检查镜像更新: %s=%s", labelKey, labelValue)

	// 获取所有带有指定标签的容器
	containers, err := c.containerSvc.GetByLabel(ctx, labelKey, labelValue)
	if err != nil {
		return nil, fmt.Errorf("获取标签容器失败: %w", err)
	}

	// 使用通用检查逻辑
	return c.checkImages(ctx, containers, utils.CreateCheckCallback())
}

// CheckAll 检查所有容器的镜像更新
func (c *Checker) CheckAll(ctx context.Context) (*types.BatchCheckResult, error) {
	logger.Info("开始检查所有容器的镜像更新")

	// 获取所有容器
	containers, err := c.containerSvc.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取所有容器失败: %w", err)
	}

	// 使用通用检查逻辑
	return c.checkImages(ctx, containers, utils.CreateCheckCallback())
}

// checkImages 通用的镜像检查逻辑
func (c *Checker) checkImages(ctx context.Context, containers []types.ContainerInfo, callback types.CheckCallback) (*types.BatchCheckResult, error) {
	startTime := time.Now()
	result := &types.BatchCheckResult{
		Containers: containers,
	}
	result.Summary.TotalContainers = len(containers)

	if len(containers) == 0 {
		logger.Warn("未找到匹配的容器")
		return result, nil
	}

	logger.Info("找到 %d 个容器，开始检查镜像更新", len(containers))

	// 提取唯一的镜像名称
	imageNames := c.containerSvc.ExtractUniqueImages(containers)
	result.Summary.TotalImages = len(imageNames)
	logger.Debug("提取到 %d 个唯一镜像: %v", len(imageNames), imageNames)

	// 并发检查所有镜像
	var wg sync.WaitGroup
	resultsChan := make(chan *types.ImageCheckResult, len(imageNames))
	errChan := make(chan error, len(imageNames))

	logger.Debug("开始并发检查 %d 个镜像", len(imageNames))

	for _, imageName := range imageNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			logger.Info("开始检查镜像: %s", name)
			info, err := c.imageSvc.CheckUpdate(ctx, name)
			if err != nil {
				logger.Debug("检查镜像 %s 失败: %v", name, err)
				errChan <- fmt.Errorf("检查镜像 %s 失败: %w", name, err)
				resultsChan <- info
				return
			}
			logger.Debug("镜像 %s 检查完成，是否有更新: %v", name, info.IsUpdated)
			resultsChan <- info
		}(imageName)
	}

	// 启动一个goroutine来收集结果并调用回调
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	// 实时收集结果并调用回调
	for info := range resultsChan {
		result.Images = append(result.Images, info)
		// 如果有回调函数，立即调用
		if callback != nil {
			callback(info)
		}
	}

	// 处理错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// 生成统计信息
	result.Summary.Duration = time.Since(startTime)

	for _, info := range result.Images {
		if info.Error != "" {
			result.Summary.Failed++
		} else if info.IsUpdated {
			result.Summary.Updated++
		} else {
			result.Summary.UpToDate++
		}
	}

	// 记录检查结果
	logger.Info("镜像检查完成: 更新 %d, 最新 %d, 失败 %d, 耗时 %v",
		result.Summary.Updated, result.Summary.UpToDate, result.Summary.Failed, result.Summary.Duration)

	// 如果有错误，返回第一个错误
	if len(errors) > 0 {
		logger.Warn("检查过程中出现 %d 个错误", len(errors))
		return result, errors[0]
	}

	return result, nil
}

// Close 关闭所有资源
func (c *Checker) Close() error {
	var errors []error

	if c.clientManager != nil {
		if err := c.clientManager.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("关闭资源时出现错误: %v", errors)
	}

	return nil
}
